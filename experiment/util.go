package main

import (
	proto "github.com/huin/mqtt"
	"github.com/tealeg/xlsx/v3"
	"log"
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
	sheetName := "Sheet1"
	sh, ok := wb.Sheet[sheetName]
	if !ok {
		sh, err := wb.AddSheet(sheetName)
		if err != nil {
			log.Fatal(err)
		}
		row := sh.AddRow()
		for _, title := range header {
			row.AddCell().SetString(title)
		}
		return sh
	}
	return sh
}
