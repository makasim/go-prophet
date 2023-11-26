package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/makasim/go-prophet"
)

func main() {
	df, err := csvToDataFrame("example/example.csv")
	if err != nil {
		panic(err)
	}

	p := prophet.New()

	fs, err := p.Forecast(df)
	if err != nil {
		panic(err)
	}

	log.Println(fs)
}

func csvToDataFrame(csvFile string) (prophet.DataFrame, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("os: open: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv: reader: read all: %w", err)
	}

	df := prophet.DataFrame{}
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) != 2 {
			return nil, fmt.Errorf("invalid record: %v", record)
		}

		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("parse float: %w", err)
		}

		df = append(df, prophet.DataPoint{
			Ds: record[0],
			Y:  y,
		})
	}

	return df, nil
}
