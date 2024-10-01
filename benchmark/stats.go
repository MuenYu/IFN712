package main

import (
	"fmt"
	"github.com/eclesh/welford"
	"github.com/tealeg/xlsx/v3"
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

func (info *testRecord) write2Row(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	row.AddCell().Value = info.proto
	row.AddCell().SetInt64(info.latency.Nanoseconds())
	row.AddCell().Value = info.errMsg
}

type statsData struct {
	successCount int
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
	messageCapacity = pairs * messagePerPair
	statsChan       = make(chan testRecord, messageCapacity)

	data = map[string]*statsData{
		"TCP": {},
		"KCP": {},
	}
	file     = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet1")
)

func runStats() {
	header := sheet.AddRow()
	header.AddCell().Value = "proto"
	header.AddCell().Value = "latency (nano seconds)"
	header.AddCell().Value = "error message"
	for record := range statsChan {
		protoStats := data[record.proto]
		if record.errMsg == "" {
			protoStats.successCount++
		}
		record.write2Row(sheet)
		protoStats.stats.Add(float64(record.latency))
	}
}

func stopStats() {
	close(statsChan)
	//date := time.Now().Format("20060102-150405")
	//fileName := fmt.Sprintf("%s-%d-%d.xlsx", date, pairs, messagePerPair)
	for proto, protoData := range data {
		log.Printf("%s:\n", proto)
		protoData.report()
	}
	//if err := file.Save(fileName); err != nil {
	//	log.Printf("error saving stats: %s\n", err)
	//}
	//log.Println("stats saved to", fileName)
}
