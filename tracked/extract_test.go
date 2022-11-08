package tracked

import (
	"encoding/csv"
	"log"
	"os"
	"testing"
)

var testDataWithoutHeader = [][]string{
	{"AAPL", "Apple Inc."},
	{"MSFT", "Microsoft Corporation"},
	{"TSLA", "Tesla Inc."},
}

var testDataWithHeader = [][]string{
	{"Symbol", "Name"},
	{"AAPL", "Apple Inc."},
	{"MSFT", "Microsoft Corporation"},
	{"TSLA", "Tesla Inc."},
}

func createTempCSVFile(dataSource [][]string) (string, error) {
	wd, _ := os.Getwd()
	tempFile, tempFileErr := os.CreateTemp(wd, "tracked_test_*.csv")

	if tempFileErr != nil {
		return "", tempFileErr
	}

	defer func(tempFile *os.File) {
		_ = tempFile.Close()
	}(tempFile)

	w := csv.NewWriter(tempFile)
	defer w.Flush()

	for _, record := range dataSource {
		if err := w.Write(record); err != nil {
			log.Fatal(err)
		}
	}

	return tempFile.Name(), nil
}

func TestExtractWithoutHeader(t *testing.T) {
	tempFileName, tempFileErr := createTempCSVFile(testDataWithoutHeader)

	if tempFileErr != nil {
		t.Fatal(tempFileErr)
	}

	defer func() {
		_ = os.Remove(tempFileName)
	}()

	actual, extractErr := Extract(tempFileName, false)

	if extractErr != nil {
		t.Fatal(extractErr)
	}

	if actual == nil {
		t.Errorf("expected: %v, actual: %v", []string{}, actual)
	}

	if len(actual) != len(testDataWithoutHeader) {
		t.Errorf("expected: %v, actual: %v", len(testDataWithoutHeader), len(actual))
	}
}

func TestExtractWithHeader(t *testing.T) {
	tempFileName, tempFileErr := createTempCSVFile(testDataWithHeader)

	if tempFileErr != nil {
		t.Fatal(tempFileErr)
	}

	defer func() {
		_ = os.Remove(tempFileName)
	}()

	actual, extractErr := Extract(tempFileName, true)

	if extractErr != nil {
		t.Fatal(extractErr)
	}

	if actual == nil {
		t.Errorf("expected: %v, actual: %v", []string{}, actual)
	}

	if len(actual) != len(testDataWithoutHeader) {
		t.Errorf("expected: %v, actual: %v", len(testDataWithoutHeader), len(actual))
	}
}
