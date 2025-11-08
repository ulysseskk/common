package csvUtil

import (
	"github.com/gocarina/gocsv"
)

func CsvFileInfoStructs(content string, datas interface{}) error {
	if err := gocsv.UnmarshalString(content, datas); err != nil { // Load clients from file
		return err
	}
	return nil
}

func StructsToCsvFileContent(datas interface{}) (string, error) {
	return gocsv.MarshalString(datas) // Get all clients as CSV string
}
