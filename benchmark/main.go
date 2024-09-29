package main

import (
	"log"
	"sync"
)

const (
	host         = "127.0.0.1:1883"
	pairs        = 100
	mode         = "KCP"
	messageCount = 100
	payloadSize  = 1024
)

var connectionsReady sync.WaitGroup

func main() {
	connectionsReady.Add(2 * pairs)
	go func() {
		connectionsReady.Wait()
		log.Println("all pub-sub pairs created")
	}()

	var wg sync.WaitGroup
	for i := 1; i <= pairs; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			switch mode {
			case "TCP":
				pingTcp(i)
			case "KCP":
				pingKcp(i)
			default:
				log.Fatalln("mode should be TCP or KCP only")
			}
		}(i)
	}
	wg.Wait()
	log.Println("all pub-sub pairs finished")
}
