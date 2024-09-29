package main

import (
	"bytes"
	"fmt"
	proto "github.com/huin/mqtt"
	"log"
	"mqtt/kcp"
	"os"
	"time"
)

func pingKcp(i int) {
	topic := fmt.Sprintf("pingtest/%v/request", i)
	topic2 := fmt.Sprintf("pingtest/%v/reply", i)

	start := make(chan struct{})
	stop := make(chan struct{})
	payload := make(proto.BytesPayload, payloadSize)

	go func() {
		cc := connectKcp()
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
	cc := connectKcp()
	defer cc.Disconnect()

	cc.Subscribe([]proto.TopicQos{
		{topic2, proto.QosAtMostOnce},
	})

	for count := 0; count < messageCount; count++ {
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
		elapsed := time.Since(timeStart)
		log.Printf("%d: %v\n", i, elapsed)

		buf := &bytes.Buffer{}
		err := in.Payload.WritePayload(buf)
		if err != nil {
			// TODO: exception data collection
			log.Printf("record error: %v\n", err)
		} else if !bytes.Equal(buf.Bytes(), payload) {
			// TODO: exception data collection
			log.Println("payload mismatch")
		} else {
			// TODO: latency data collection
		}
	}
	close(stop)
}

func connectKcp() *kcp.ClientConn {
	conn, err := kcp.DialKCP(host)
	if err != nil {
		fmt.Printf("kcp dial: %v\n", err)
		os.Exit(2)
	}
	cc := kcp.NewClientConn(conn)
	err = cc.Connect("", "")
	if err != nil {
		fmt.Printf("kcp connect: %v\n", err)
		os.Exit(3)
	}
	connectionsReady.Done()
	return cc
}
