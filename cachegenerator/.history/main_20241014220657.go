package main

import (
	"encoding/csv"
	"fmt"
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
	Close []float64
	Open []float64
	Volume []float64
}


var globaldata []Record = make([]Record, 0)

func avgOfArray(array []int) float64 {
    var sum int
    for _, value := range array {
        sum += value
    }
    return float64(sum) / float64(len(array))
}

func process(year int, comc chan<- WorkerResult) {
	companies := pie.Unique(pie.Map(globaldata, func(record Record) string { return record.Ticker }))
	print("Preparing data for year ", year, "...\n")
	companydict := make(map[string][]Record)

	for i,row := range globaldata {
		if companydict[row.Ticker] == nil {
			companydict[row.Ticker] = make([]Record, 0)
		} else {
			companydict[row.Ticker] = append(companydict[row.Ticker], globaldata[i])
		}
	}

	close := make([]float64, len(companies))
	open := make([]float64, len(companies))
	volume := make([]float64, len(companies))
	print("Processing data for year ", year, "...\n")
	for i, company := range companies {
		fetch := pie.Filter(companydict[company], func(record Record) bool {
			return record.Date.Year() == year
		})
		close[i] = avgOfArray(pie.Map(fetch, func(record Record) int { return int(record.Close) }))
		open[i] = avgOfArray(pie.Map(fetch, func(record Record) int { return int(record.Open) }))
		volume[i] = avgOfArray(pie.Map(fetch, func(record Record) int { return record.Volume }))
	}
	comc <- WorkerResult{
		Year: year,
		Close: close,
		Open: open,
		Volume: volume,
	}

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


	years := pie.Unique(
		pie.Map(converted, func(record Record) int {
			return record.Date.Year()
		}),
	)
	fmt.Println("Processing data for years: ", years)

	results := make(chan WorkerResult, len(years))
	for _, year := range(years) {
		go process(year, results)
	}
	resultsdata := make([]WorkerResult, len(years))
	for i := 0; i < len(years); i++ {
		resultsdata[i] = <-results
	}
}