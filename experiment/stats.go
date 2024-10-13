package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/tealeg/xlsx/v3"
	"log"
	"sync"
	"time"
)

// one request-response trip record: pub -> broker -> sub -> broker -> pub
type testRecord struct {
	proto   string
	latency time.Duration
	errMsg  string
}

func (record *testRecord) String() string {
	return fmt.Sprintf("proto: %s, latency: %v, errMsg: %s", record.proto, record.latency, record.errMsg)
}

func (record *testRecord) write2Row(sheet *xlsx.Sheet) {
	row := sheet.AddRow()
	row.AddCell().Value = record.proto
	row.AddCell().SetInt64(int64(record.latency))
	row.AddCell().Value = record.errMsg
	row.AddCell().SetInt64(int64(*payloadSize))
	row.AddCell().SetInt64(int64(*reqInterval / 1e6))
	row.AddCell().Value = *network
}

var (
	recordChan = make(chan testRecord, 2**messagePerPair)
	wg         = new(sync.WaitGroup)
	file       = openOrCreateXlsx()
	sheet      = initSheetAndHead(file, []string{
		"protocol",
		"latency (nanosecond)",
		"error message",
		"payload",
		"interval",
		"network",
	})
)

func runStats() {
	bar := progressbar.Default(int64(2 * *messagePerPair))
	for record := range recordChan {
		wg.Add(1)
		go func(record testRecord) {
			defer wg.Done()
			defer bar.Add(1)
			record.write2Row(sheet)
		}(record)
	}
}

func outputReport() {
	wg.Wait()
	close(recordChan)
	if err := file.Save(*outputFile); err != nil {
		log.Println(err.Error())
	}
}
