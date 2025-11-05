package main

import (
	"encoding/csv"
	"os"
)

func main() {
	// read the csv file
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	data = csv.NewReader(file)
	records, err := data.ReadAll()

}