package main

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/elliotchance/pie/v2"
)

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
	companydict 
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
	parseddate := []time.Time{}
	for i, date := range getColumn(records, "Date") {
		parseddate[i], err = time.Parse("2006-01-02", date)
		if err != nil {
			parseddate[i] = time.Time{}
		}
	}

	companies := pie.Unique(getColumn(records, "Ticker"))
	years := pie.Unique(
		pie.Map(parseddate, func(date time.Time) int {
			return date.Year()
		}),
	)




}