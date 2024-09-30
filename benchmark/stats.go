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

func (sd *statsData) write2Xlsx(sheet *xlsx.Sheet) {
	for _, record := range sd.recordList {
		row := sheet.AddRow()
		row.AddCell().Value = record.proto
		row.AddCell().SetInt64(record.latency.Nanoseconds())
		row.AddCell().Value = record.errMsg
	}
}

var (
	messageCapacity = pairs * messagePerPair
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
	date := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("%s-%d-%d.xlsx", date, pairs, messagePerPair)
	file, sheet := prepSheet()
	for proto, protoData := range data {
		log.Printf("%s:\n", proto)
		protoData.report()
		protoData.write2Xlsx(sheet)
	}
	if err := file.Save(fileName); err != nil {
		log.Printf("error saving stats: %s\n", err)
	}
	log.Println("stats saved to", fileName)

}

func prepSheet() (*xlsx.File, *xlsx.Sheet) {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := sheet.AddRow()
	header.AddCell().Value = "proto"
	header.AddCell().Value = "latency (nano seconds)"
	header.AddCell().Value = "error message"
	return file, sheet
}
