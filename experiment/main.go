package main

import (
	"time"
)

const (
	//host = "127.0.0.1:1883"
	host           = "3.24.217.226:1883"
	messagePerPair = 10
	payloadSize    = 60
	timeout        = time.Second
	reqInterval    = 100 * time.Millisecond

	outputFile = "data.xlsx"
)

func main() {
	go runStats()
	pingTcp()
	pingKcp()
	outputReport()
}
