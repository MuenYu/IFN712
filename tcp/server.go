package tcp

import (
	"fmt"
	proto "github.com/huin/mqtt"
	"io"
	"log"
	"math/rand"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

type subscriptions struct {
	workers int
	posts   chan post

	mu        sync.Mutex // guards access to fields below
	subs      map[string][]*incomingConn
	wildcards []wild
	retain    map[string]retain
	stats     *stats
}

func newSubscriptions(workers int) *subscriptions {
	s := &subscriptions{
		subs:    make(map[string][]*incomingConn),
		retain:  make(map[string]retain),
		posts:   make(chan post, postQueue),
		workers: workers,
	}
	for i := 0; i < s.workers; i++ {
		go s.run(i)
	}
	return s
}

func (s *subscriptions) sendRetain(topic string, c *incomingConn) {
	s.mu.Lock()
	var tlist []string
	if isWildcard(topic) {

		// TODO: select matching topics from the retain map
	} else {
		tlist = []string{topic}
	}
	for _, t := range tlist {
		if r, ok := s.retain[t]; ok {
			c.submit(&r.m)
		}
	}
	s.mu.Unlock()
}

func (s *subscriptions) add(topic string, c *incomingConn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if isWildcard(topic) {
		w := newWild(topic, c)
		if w.valid() {
			s.wildcards = append(s.wildcards, w)
		}
	} else {
		s.subs[topic] = append(s.subs[topic], c)
	}
}

// Find all connections that are subscribed to this topic.
func (s *subscriptions) subscribers(topic string) []*incomingConn {
	s.mu.Lock()
	defer s.mu.Unlock()

	// non-wildcard subscribers
	res := s.subs[topic]

	// process wildcards
	parts := strings.Split(topic, "/")
	for _, w := range s.wildcards {
		if w.matches(parts) {
			res = append(res, w.c)
		}
	}

	return res
}

// Remove all subscriptions that refer to a connection.
func (s *subscriptions) unsubAll(c *incomingConn) {
	s.mu.Lock()
	for _, v := range s.subs {
		for i := range v {
			if v[i] == c {
				v[i] = nil
			}
		}
	}

	// remove any associated entries in the wildcard list
	var wildNew []wild
	for i := 0; i < len(s.wildcards); i++ {
		if s.wildcards[i].c != c {
			wildNew = append(wildNew, s.wildcards[i])
		}
	}
	s.wildcards = wildNew

	s.mu.Unlock()
}

// Remove the subscription to topic for a given connection.
func (s *subscriptions) unsub(topic string, c *incomingConn) {
	s.mu.Lock()
	if subs, ok := s.subs[topic]; ok {
		nils := 0

		// Search the list, removing references to our connection.
		// At the same time, count the nils to see if this list is now empty.
		for i := 0; i < len(subs); i++ {
			if subs[i] == c {
				subs[i] = nil
			}
			if subs[i] == nil {
				nils++
			}
		}

		if nils == len(subs) {
			delete(s.subs, topic)
		}
	}
	s.mu.Unlock()
}

// The subscription processing worker.
func (s *subscriptions) run(id int) {
	tag := fmt.Sprintf("tcp worker %d ", id)
	log.Print(tag, "started")
	for post := range s.posts {
		// Remember the original retain setting, but send out immediate
		// copies without retain: "When a server sends a PUBLISH to a client
		// as a result of a subscription that already existed when the
		// original PUBLISH arrived, the Retain flag should not be set,
		// regardless of the Retain flag of the original PUBLISH.
		isRetain := post.m.Header.Retain
		post.m.Header.Retain = false

		// Handle "retain with payload size zero = delete retain".
		// Once the delete is done, return instead of continuing.
		if isRetain && post.m.Payload.Size() == 0 {
			s.mu.Lock()
			delete(s.retain, post.m.TopicName)
			s.mu.Unlock()
			return
		}

		// Find all the connections that should be notified of this message.
		conns := s.subscribers(post.m.TopicName)

		// Queue the outgoing messages
		for _, c := range conns {
			// Do not echo messages back to where they came from.
			if c == post.c {
				continue
			}

			if c != nil {
				c.submit(post.m)
			}
		}

		if isRetain {
			s.mu.Lock()
			// Save a copy of it, and set that copy's Retain to true, so that
			// when we send it out later we notify new subscribers that this
			// is an old message.
			msg := *post.m
			msg.Header.Retain = true
			s.retain[post.m.TopicName] = retain{m: msg}
			s.mu.Unlock()
		}
	}
}

func (s *subscriptions) submit(c *incomingConn, m *proto.Publish) {
	s.posts <- post{c: c, m: m}
}

// A post is a unit of work for the subscription processing workers.
type post struct {
	c *incomingConn
	m *proto.Publish
}

// A Server holds all the state associated with an MQTT server.
type Server struct {
	l             net.Listener
	subs          *subscriptions
	stats         *stats
	Done          chan struct{}
	StatsInterval time.Duration // Defaults to 10 seconds. Must be set using sync/atomic.StoreInt64().
	Dump          bool          // When true, dump the messages in and out.
	rand          *rand.Rand
}

// NewServer creates a new MQTT server, which accepts connections from
// the given listener. When the server is stopped (for instance by
// another goroutine closing the net.Listener), channel Done will become
// readable.
func NewServer(l net.Listener) *Server {
	svr := &Server{
		l:             l,
		stats:         &stats{},
		Done:          make(chan struct{}),
		StatsInterval: time.Second * 10,
		subs:          newSubscriptions(runtime.GOMAXPROCS(0)),
	}

	// start the stats reporting goroutine
	go func() {
		for {
			svr.stats.publish(svr.subs, svr.StatsInterval)
			select {
			case <-svr.Done:
				return
			default:
				// keep going
			}
			time.Sleep(svr.StatsInterval)
		}
	}()

	return svr
}

// Start makes the Server start accepting and handling connections.
func (s *Server) Start() {
	go func() {
		defer close(s.Done)
		for {
			conn, err := s.l.Accept()
			if err != nil {
				log.Print("Accept: ", err)
			}

			cli := s.newIncomingConn(conn)
			s.stats.clientConnect()
			cli.start()
		}
	}()
}

// An IncomingConn represents a connection into a Server.
type incomingConn struct {
	svr      *Server
	conn     net.Conn
	jobs     chan job
	clientid string
	Done     chan struct{}
}

var clients = make(map[string]*incomingConn)
var clientsMu sync.Mutex

// newIncomingConn creates a new incomingConn associated with this
// server. The connection becomes the property of the incomingConn
// and should not be touched again by the caller until the Done
// channel becomes readable.
func (s *Server) newIncomingConn(conn net.Conn) *incomingConn {
	return &incomingConn{
		svr:  s,
		conn: conn,
		jobs: make(chan job, sendingQueueLength),
		Done: make(chan struct{}),
	}
}

// Start reading and writing on this connection.
func (c *incomingConn) start() {
	go c.reader()
	go c.writer()
}

// Add this connection to the map, or find out that an existing connection
// already exists for the same client-id.
func (c *incomingConn) add() *incomingConn {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	existing, ok := clients[c.clientid]
	if ok {
		// this client id already exists, return it
		return existing
	}

	clients[c.clientid] = c
	return nil
}

// Delete a connection; the connection must be closed by the caller first.
func (c *incomingConn) del() {
	clientsMu.Lock()
	delete(clients, c.clientid)
	clientsMu.Unlock()
	return
}

// Queue a message; no notification of sending is done.
func (c *incomingConn) submit(m proto.Message) {
	j := job{m: m}
	select {
	case c.jobs <- j:
	default:
		log.Print(c, ": failed to submit message")
	}
	return
}

func (c *incomingConn) String() string {
	return fmt.Sprintf("{IncomingConn: %v}", c.clientid)
}

// Queue a message, returns a channel that will be readable
// when the message is sent.
func (c *incomingConn) submitSync(m proto.Message) receipt {
	j := job{m: m, r: make(receipt)}
	c.jobs <- j
	return j.r
}

func (c *incomingConn) reader() {
	// On exit, close the connection and arrange for the writer to exit
	// by closing the output channel.
	defer func() {
		c.conn.Close()
		c.svr.stats.clientDisconnect()
		close(c.jobs)
	}()

	for {
		// TODO: timeout (first message and/or keepalives)
		m, err := proto.DecodeOneMessage(c.conn, nil)
		if err != nil {
			if err == io.EOF {
				return
			}
			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				return
			}
			log.Print("reader: ", err)
			return
		}
		c.svr.stats.messageRecv()

		if c.svr.Dump {
			log.Printf("dump  in: %T", m)
		}

		switch m := m.(type) {
		case *proto.Connect:
			rc := proto.RetCodeAccepted

			if m.ProtocolName != "MQIsdp" ||
				m.ProtocolVersion != 3 {
				log.Print("reader: reject connection from ", m.ProtocolName, " version ", m.ProtocolVersion)
				rc = proto.RetCodeUnacceptableProtocolVersion
			}

			// Check client id.
			if len(m.ClientId) < 1 || len(m.ClientId) > 23 {
				rc = proto.RetCodeIdentifierRejected
			}
			c.clientid = m.ClientId

			// Disconnect existing connections.
			if existing := c.add(); existing != nil {
				disconnect := &proto.Disconnect{}
				r := existing.submitSync(disconnect)
				r.wait()
				c.add()
			}

			// TODO: Last will

			connack := &proto.ConnAck{
				ReturnCode: rc,
			}
			c.submit(connack)

			// close connection if it was a bad connect
			if rc != proto.RetCodeAccepted {
				log.Printf("Connection refused for %v: %v", c.conn.RemoteAddr(), ConnectionErrors[rc])
				return
			}

			// Log in mosquitto format.
			clean := 0
			if m.CleanSession {
				clean = 1
			}
			log.Printf("New client connected from %v as %v (c%v, k%v).", c.conn.RemoteAddr(), c.clientid, clean, m.KeepAliveTimer)

		case *proto.Publish:
			// TODO: Proper QoS support
			if m.Header.QosLevel != proto.QosAtMostOnce {
				log.Printf("reader: no support for QoS %v yet", m.Header.QosLevel)
				return
			}
			if m.Header.QosLevel != proto.QosAtMostOnce && m.MessageId == 0 {
				// Invalid message ID. See MQTT-2.3.1-1.
				log.Printf("reader: invalid MessageId in PUBLISH.")
				return
			}
			if isWildcard(m.TopicName) {
				log.Print("reader: ignoring PUBLISH with wildcard topic ", m.TopicName)
			} else {
				c.svr.subs.submit(c, m)
			}
			c.submit(&proto.PubAck{MessageId: m.MessageId})

		case *proto.PingReq:
			c.submit(&proto.PingResp{})

		case *proto.Subscribe:
			if m.Header.QosLevel != proto.QosAtLeastOnce {
				// protocol error, disconnect
				return
			}
			if m.MessageId == 0 {
				// Invalid message ID. See MQTT-2.3.1-1.
				log.Printf("reader: invalid MessageId in SUBSCRIBE.")
				return
			}
			suback := &proto.SubAck{
				MessageId: m.MessageId,
				TopicsQos: make([]proto.QosLevel, len(m.Topics)),
			}
			for i, tq := range m.Topics {
				// TODO: Handle varying QoS correctly
				c.svr.subs.add(tq.Topic, c)
				suback.TopicsQos[i] = proto.QosAtMostOnce
			}
			c.submit(suback)

			// Process retained messages.
			for _, tq := range m.Topics {
				c.svr.subs.sendRetain(tq.Topic, c)
			}

		case *proto.Unsubscribe:
			if m.Header.QosLevel != proto.QosAtMostOnce && m.MessageId == 0 {
				// Invalid message ID. See MQTT-2.3.1-1.
				log.Printf("reader: invalid MessageId in UNSUBSCRIBE.")
				return
			}
			for _, t := range m.Topics {
				c.svr.subs.unsub(t, c)
			}
			ack := &proto.UnsubAck{MessageId: m.MessageId}
			c.submit(ack)

		case *proto.Disconnect:
			return

		default:
			log.Printf("reader: unknown msg type %T", m)
			return
		}
	}
}

func (c *incomingConn) writer() {

	// Close connection on exit in order to cause reader to exit.
	defer func() {
		c.conn.Close()
		c.del()
		c.svr.subs.unsubAll(c)
	}()

	for job := range c.jobs {
		if c.svr.Dump {
			log.Printf("dump out: %T", job.m)
		}

		// TODO: write timeout
		err := job.m.Encode(c.conn)
		if job.r != nil {
			// notifiy the sender that this message is sent
			close(job.r)
		}
		if err != nil {
			// This one is not interesting; it happens when clients
			// disappear before we send their acks.
			oe, isoe := err.(*net.OpError)
			if isoe && oe.Err.Error() == "use of closed network connection" {
				return
			}
			// In Go < 1.5, the error is not an OpError.
			if err.Error() == "use of closed network connection" {
				return
			}

			log.Print("writer: ", err)
			return
		}
		c.svr.stats.messageSend()

		if _, ok := job.m.(*proto.Disconnect); ok {
			log.Print("writer: sent disconnect message")
			return
		}
	}
}
