package main

import (
	"bytes"
	"fmt"
	proto "github.com/huin/mqtt"
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
		elapsed := time.Since(timeStart)

		buf := &bytes.Buffer{}
		err := in.Payload.WritePayload(buf)
		if err != nil {
			// TODO: exception data collection
			statsChan <- testRecord{
				proto:   "KCP",
				latency: elapsed,
				errMsg:  err.Error(),
			}
		} else if !bytes.Equal(buf.Bytes(), payload) {
			// TODO: exception data collection
			statsChan <- testRecord{
				proto:   "KCP",
				latency: elapsed,
				errMsg:  "payload mismatch",
			}
		} else {
			// TODO: latency data collection
			statsChan <- testRecord{
				proto:   "KCP",
				latency: elapsed,
				errMsg:  "",
			}
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
	// enable kcp turbo mode
	conn.SetNoDelay(1, 10, 2, 1)
	cc := kcp.NewClientConn(conn)
	err = cc.Connect("", "")
	if err != nil {
		fmt.Printf("kcp connect: %v\n", err)
		os.Exit(3)
	}
	connectionsReady.Done()
	return cc
}
