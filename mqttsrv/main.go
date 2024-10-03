package main

import (
	"flag"
	"log"
	"mqtt/kcp"
	"mqtt/tcp"
	"net"
)

var addr = flag.String("addr", ":1883", "listen address of broker")

func main() {
	flag.Parse()

	lTcp, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Print("listen: ", err)
		return
	}
	lKcp, err := kcp.ListenKCP(*addr)
	if err != nil {
		log.Print("listen: ", err)
		return
	}
	svrTcp := tcp.NewServer(lTcp)
	svrKcp := kcp.NewServer(lKcp)

	svrTcp.Start()
	svrKcp.Start()

	select {
	case <-svrTcp.Done:
	case <-svrKcp.Done:
		log.Println("server is closed")
	}
}
