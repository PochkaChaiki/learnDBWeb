package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	bookname := "test.xlsx"
	sheetname := "Sheet1"

	f, err := excelize.OpenFile(bookname)
	if err != nil {
		log.Fatalf("excel open error: %s", err)
		return

	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("excel close error: %s", err)
			return
		}
	}()

	rows, err := f.GetRows(sheetname)
	if err != nil {
		log.Fatalf("excel read error: %s", err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Printf(colCell, "\t")
		}
		fmt.Println()
	}

}
