package main

import (
	"github.com/schollz/progressbar/v3"
	"log"
	"sync"
)

const (
	host           = "127.0.0.1:1883"
	pairs          = 100
	messagePerPair = 100
	payloadSize    = 1024
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
		bar := progressbar.Default(pairs)
		for id := 1; id <= pairs; id++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				ping(id)
				if err := bar.Add(1); err != nil {
					log.Println(err.Error())
				}
			}(id)
		}
		wg.Wait()
	}
	stopStats()
}
