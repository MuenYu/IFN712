package main

import (
	"log"
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
	testSet          = map[string]func(i int){
		"TCP": pingTcp,
		"KCP": pingKcp,
	}
)

func main() {
	for proto, pingFunc := range testSet {
		log.Printf("-------------%s start-----------------\n", proto)
		connectionsReady.Add(2 * pairs)
		go func() {
			connectionsReady.Wait()
			log.Println("all pub-sub pairs created")
		}()
		var wg sync.WaitGroup
		for id := 1; id <= pairs; id++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				pingFunc(id)
			}(id)
		}
		wg.Wait()
		log.Printf("-------------%s end-----------------\n", proto)
	}
}
