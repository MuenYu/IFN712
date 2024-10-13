package main

import (
	"flag"
	"time"
)

var (
	host           = flag.String("host", "127.0.0.1:1883", "mqtt broker host")
	messagePerPair = flag.Int("messages", 100, "the quantity of requests going to send for each")
	payloadSize    = flag.Int("payload", 60, "the payload size for requests")
	timeout        = flag.Duration("timeout", time.Second, "the timeout for each request")
	reqInterval    = flag.Duration("interval", 100*time.Millisecond, "the interval between requests")
	outputFile     = flag.String("output", "data.xlsx", "the file to store testing result")
	network        = flag.String("network", "ethernet", "e.g. ethernet, wifi")
)

func main() {
	flag.Parse()
	go runStats()
	pingTcp()
	pingKcp()
	outputReport()
}
