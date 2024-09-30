package main

import (
	"github.com/schollz/progressbar/v3"
	"sync"
)

const (
	host         = "127.0.0.1:1883"
	pairs        = 100
	messageCount = 100
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
		bar := progressbar.Default(pairs)
		var wg sync.WaitGroup
		for id := 1; id <= pairs; id++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				defer bar.Add(1)
				ping(id)
			}(id)
		}
		wg.Wait()
	}
	stopStats()
}
