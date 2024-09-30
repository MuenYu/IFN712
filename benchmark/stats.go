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
	successCount int
	recordList   []testRecord
	stats        welford.Stats
}

func (sd *statsData) report() {
	log.Printf("success rate: %d/%d\n", sd.successCount, messageCapacity)
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
			recordList: make([]testRecord, 0, messageCapacity),
		},
		"KCP": {
			recordList: make([]testRecord, 0, messageCapacity),
		},
	}
)

func runStats() {
	for record := range statsChan {
		protoStats := data[record.proto]
		protoStats.recordList = append(protoStats.recordList, record)
		if record.errMsg == "" {
			protoStats.successCount++
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
