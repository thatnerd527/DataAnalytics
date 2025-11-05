package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"runtime"
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

type WorkerRange struct {
	Min int64
	Max int64
}

func workerSplitter(workers int, amount int64) []WorkerRange {
	amounts := make([]int64, workers)
	for i := int64(0); i < amount ; i++ {
		amounts[i % int64(workers)]++
	}
	ranges := make([]WorkerRange, workers)
	var min int64
	for i, amount := range amounts {
		ranges[i] = WorkerRange{
			Min: min,
			Max: min + amount,
		}
		min += amount
	}
	return ranges
}

var globaldata []Record = make([]Record, 0)
var companydict map[string][]Record = make(map[string][]Record)

func avgOfArray(array []float64) float64 {
    var sum float64
    for _, value := range array {
        sum += value
    }
    return float64(sum) / float64(len(array))
}

func process(year int, comc chan<- WorkerResult) {
	companies := pie.Unique(pie.Map(globaldata, func(record Record) string { return record.Ticker }))
	print("Preparing data for year ", year, "...\n")


	close := make([]float64, len(companies))
	open := make([]float64, len(companies))
	volume := make([]float64, len(companies))
	print("Processing data for year ", year, "...\n")
	for i, company := range companies {
		fetch := pie.Filter(companydict[company], func(record Record) bool {
			return record.Date.Year() == year
		})
		if len(fetch) == 0 {
			close[i] = 0
			open[i] = 0
			volume[i] = 0
			continue
		}
		avg1 := avgOfArray(pie.Map(fetch, func(record Record) float64 { return record.Close }))
		close[i] = avg1
		open[i] = avgOfArray(pie.Map(fetch, func(record Record) float64 { return record.Open }))
		volume[i] = avgOfArray(pie.Map(fetch, func(record Record) float64 { return float64(record.Volume) }))
	}
	print("Data for year ", year, " processed\n")
	comc <- WorkerResult{
		Year: year,
		Close: close,
		Open: open,
		Volume: volume,
	}

}

func main() {
	// Part 1
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
	for i,row := range globaldata {
		if companydict[row.Ticker] == nil {
			companydict[row.Ticker] = make([]Record, 0)
		} else {
			companydict[row.Ticker] = append(companydict[row.Ticker], globaldata[i])
		}
	}
	fmt.Println("Processing data for years: ", years)

	results := make(chan WorkerResult, len(years))
	for _, year := range(years) {
		go process(year, results)
	}
	resultsdata := make([]WorkerResult, len(years))
	for i := 0; i < len(years); i++ {
		resultsdata[i] = <-results
	}
	peryearclose := map[int][]float64{}
	peryearopen := map[int][]float64{}
	peryearvolume := map[int][]float64{}

	for _, result := range resultsdata {
		peryearclose[result.Year] = result.Close
		peryearopen[result.Year] = result.Open
		peryearvolume[result.Year] = result.Volume
	}
	marshalled, err := json.Marshal(peryearclose)
	if (err != nil) {
		panic(err)
	}
	os.WriteFile("peryearopen.json", marshalled, 0644)
	marshalled, err = json.Marshal(peryearopen)
	if (err != nil) {
		panic(err)
	}
	os.WriteFile("peryearclose.json", marshalled, 0644)
	marshalled, err = json.Marshal(peryearvolume)
	if (err != nil) {
		panic(err)
	}
	os.WriteFile("peryearvolume.json", marshalled, 0644)

	// 2nd part

	lowestperiod := time.Unix(0, pie.Min(pie.Map(globaldata, func(record Record) int64 { return record.Date.UnixNano() })))
	highestperiod := time.Unix(0, pie.Max(pie.Map(globaldata, func(record Record) int64 { return record.Date.UnixNano() })))
	daysbetween := math.Floor(highestperiod.Sub(lowestperiod).Hours() / 24)
	fmt.Println("Lowest date: ", lowestperiod)
	fmt.Println("Highest date: ", highestperiod)
	fmt.Println("Days between: ", daysbetween)
	ranges := workerSplitter(runtime.NumCPU(), int64(daysbetween))


}