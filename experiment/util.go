package main

import (
	proto "github.com/huin/mqtt"
	"github.com/tealeg/xlsx/v3"
	"log"
	"time"
)

// clearChan: clean all existing msg in channel, avoiding count duplicated requests & latency
func clearChan(incoming chan *proto.Publish) {
	for {
		select {
		case <-incoming:
		default:
			return
		}
	}
}

func openOrCreateXlsx() *xlsx.File {
	wb, err := xlsx.OpenFile(*outputFile)
	if err != nil {
		wb = xlsx.NewFile()
	}
	return wb
}

func initSheetAndHead(wb *xlsx.File, header []string) *xlsx.Sheet {
	sheetName := time.Now().Format("20060102150405")
	sheet, err := wb.AddSheet(sheetName)
	if err != nil {
		log.Fatal(err)
	}
	row := sheet.AddRow()
	for _, title := range header {
		row.AddCell().SetString(title)
	}
	return sheet
}
