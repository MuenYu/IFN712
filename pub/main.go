package main

import (
	"flag"
	"fmt"
	"mqtt/kcp"
	"os"
	"time"

	proto "github.com/huin/mqtt"
)

var host = flag.String("host", "localhost:1883", "hostname of broker")
var user = flag.String("user", "", "username")
var pass = flag.String("pass", "", "password")
var dump = flag.Bool("dump", false, "dump messages?")
var retain = flag.Bool("retain", false, "retain message?")

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: pub topic message")
		return
	}

	//conn, err := net.Dial("tcp", *host)
	conn, err := kcp.DialKCP(*host)
	if err != nil {
		fmt.Fprint(os.Stderr, "dial: ", err)
		return
	}
	//cc := tcp.NewClientConn(conn)
	cc := kcp.NewClientConn(conn)
	defer cc.Disconnect()
	cc.Dump = *dump

	if err := cc.Connect(*user, *pass); err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected with client id", cc.ClientId)

	for {
		time.Sleep(time.Second)
		msg := "hello world"
		cc.Publish(&proto.Publish{
			Header:    proto.Header{Retain: *retain},
			TopicName: flag.Arg(0),
			Payload:   proto.BytesPayload(msg),
		})
		fmt.Printf("send msg: %s\n", msg)
	}
}
