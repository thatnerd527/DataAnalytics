package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/elliotchance/pie/v2"
)

type Record struct {
	Date  time.Time
	Open  float64
	High  float64
	Low   float64
	Close float64
	Volume int
	AdjClose float64 `csv:"Adj Close"`
	Ticker string
}

type WorkerResult struct {
	Year int
	

var globaldata []Record = make([]Record, 0)

func getColumn(records [][]string, columnname string) []string {
	// get the index of the column
	index := -1
	for i, column := range records[0] {
		if column == columnname {
			index = i
			break
		}
	}
	if index == -1 {
		panic("column not found")
	}

	// get the column
	column := make([]string, len(records)-1)
	for i, record := range records[1:] {
		column[i] = record[index]
	}
	return column
}

func process(year int) {

	print("Preparing data for year ", year, "...\n")
	companydict := make(map[string][]Record)

}

func main() {
	// read the csv file
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	data := csv.NewReader(file)
	records, err := data.ReadAll()

	if err != nil {
		panic(err)
	}
	converted := pie.Map(records, func(record []string) Record {
		date, _ := time.Parse("2006-01-02", record[0])
		ticker := record[1]
		open, _ := strconv.ParseFloat(record[2], 64)
		high, _ := strconv.ParseFloat(record[3], 64)
		low, _ := strconv.ParseFloat(record[4], 64)
		close, _ := strconv.ParseFloat(record[5], 64)
		adjclose, _ := strconv.ParseFloat(record[6], 64)
		volume, _ := strconv.Atoi(record[7])
		return Record{
			Date: date,
			Open: open,
			High: high,
			Low: low,
			Close: close,
			Volume: volume,
			AdjClose: adjclose,
			Ticker: ticker,
		}
	})
	globaldata = converted

	companies := pie.Unique(pie.Map(converted, func(record Record) string { return record.Ticker }))
	years := pie.Unique(
		pie.Map(converted, func(record Record) int {
			return record.Date.Year()
		}),
	)

	for year := range(years) {
		go process(year)
	}




}