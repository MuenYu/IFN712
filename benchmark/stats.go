package main

import (
	"fmt"
	"github.com/eclesh/welford"
	"log"
	"time"
)

type testRecord struct {
	proto   string
	latency time.Duration
	errMsg  string
}

func (info *testRecord) String() string {
	return fmt.Sprintf("proto: %s, latency: %v, errMsg: %s", info.proto, info.latency, info.errMsg)
}

type statsData struct {
	successList []testRecord
	failList    []testRecord
	stats       welford.Stats
}

func (sd *statsData) report() {
	log.Printf("success rate: %.2f\n", float64(len(sd.successList))/float64(messageCapacity))
	log.Printf("min latency: %.2f\n", sd.stats.Min())
	log.Printf("max latency: %.2f\n", sd.stats.Max())
	log.Printf("mean latency: %.2f\n", sd.stats.Mean())
	log.Printf("stddev latency: %.2f\n\n", sd.stats.Stddev())
}

var (
	messageCapacity = pairs * messageCount
	statsChan       = make(chan testRecord, messageCapacity)

	data = map[string]*statsData{
		"TCP": {
			successList: make([]testRecord, 0, messageCapacity),
			failList:    make([]testRecord, 0, messageCapacity),
		},
		"KCP": {
			successList: make([]testRecord, 0, messageCapacity),
			failList:    make([]testRecord, 0, messageCapacity),
		},
	}
)

func runStats() {
	for record := range statsChan {
		protoStats := data[record.proto]
		if record.errMsg == "" {
			protoStats.successList = append(protoStats.successList, record)
		} else {
			protoStats.failList = append(protoStats.failList, record)
		}
		protoStats.stats.Add(float64(record.latency))
	}
}

func stopStats() {
	close(statsChan)
	for proto, protoData := range data {
		log.Printf("%s:\n", proto)
		protoData.report()
	}
	// TODO: write records to xlsx or csv
}
