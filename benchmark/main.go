package main

import (
	"sync"
)

const (
	host         = "127.0.0.1:1883"
	pairs        = 10
	messageCount = 1000
	payloadSize  = 1024
)

var (
	connectionsReady sync.WaitGroup
	pings            = []func(i int){
		pingTcp,
		pingKcp,
	}
)

func main() {
	go runStats()

	for _, ping := range pings {
		connectionsReady.Add(2 * pairs)
		var wg sync.WaitGroup
		for id := 1; id <= pairs; id++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				ping(id)
			}(id)
		}
		wg.Wait()
	}
	stopStats()
}
