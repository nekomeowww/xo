package xo

import (
	"encoding/csv"
	"os"
)

func ReadCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return make([][]string, 0), err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ','
	csvReader.LazyQuotes = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return make([][]string, 0), err
	}

	return records, nil
}
