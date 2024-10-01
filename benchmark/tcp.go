package main

import (
	"bytes"
	"fmt"
	proto "github.com/huin/mqtt"
	"mqtt/tcp"
	"net"
	"os"
	"time"
)

func pingTcp(i int) {
	topic := fmt.Sprintf("pingtest/%v/request", i)
	topic2 := fmt.Sprintf("pingtest/%v/reply", i)

	start := make(chan struct{})
	stop := make(chan struct{})
	payload := make(proto.BytesPayload, payloadSize)

	go func() {
		cc := connectTcp()
		defer cc.Disconnect()

		cc.Subscribe([]proto.TopicQos{
			{topic, proto.QosAtLeastOnce},
		})

		close(start)
		for {
			select {
			case in := <-cc.Incoming:
				in.TopicName = topic2
				in.Payload = payload
				cc.Publish(in)
			case <-stop:
				return
			}
		}
	}()

	<-start // wait for subscriber to be ready
	cc := connectTcp()
	defer cc.Disconnect()

	cc.Subscribe([]proto.TopicQos{
		{topic2, proto.QosAtMostOnce},
	})

	for count := 0; count < messagePerPair; count++ {
		timeStart := time.Now()
		cc.Publish(&proto.Publish{
			Header:    proto.Header{QosLevel: proto.QosAtMostOnce},
			TopicName: topic,
			Payload:   payload,
		})

		in := <-cc.Incoming
		if in == nil {
			break
		}
		record := testRecord{
			proto:   "TCP",
			latency: time.Since(timeStart),
		}

		buf := &bytes.Buffer{}
		err := in.Payload.WritePayload(buf)
		if err != nil {
			record.errMsg = err.Error()
		} else if !bytes.Equal(buf.Bytes(), payload) {
			record.errMsg = "payload mismatch"
		}
		statsChan <- record
	}
	close(stop)
}

func connectTcp() *tcp.ClientConn {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Printf("tcp dial: %v\n", err)
		os.Exit(2)
	}
	cc := tcp.NewClientConn(conn)
	err = cc.Connect("", "")
	if err != nil {
		fmt.Printf("tcp connect: %v\n", err)
		os.Exit(3)
	}
	connectionsReady.Done()
	return cc
}
