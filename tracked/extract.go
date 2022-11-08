package tracked

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// Extract - parses a CSV file where the first column contains a list of symbols, tickers or currency pairs to be tracked.
func Extract(sourceFile string, hasHeader bool) ([]string, error) {
	if _, err := os.Stat(sourceFile); err != nil {
		return nil, err
	}

	f, openErr := os.Open(sourceFile)

	if openErr != nil {
		return nil, openErr
	}

	r := csv.NewReader(f)

	symbols := make([]string, 0)

	headerSkipped := false

	for {
		record, recordErr := r.Read()

		if recordErr == io.EOF {
			break
		}

		if recordErr != nil {
			return nil, recordErr
		}

		if !headerSkipped {
			if hasHeader {
				headerSkipped = true
				continue
			}
		}

		if len(record) == 0 {
			return nil, fmt.Errorf("empty record found during extraction")
		}

		symbol := record[0]

		if symbol == "" {
			return nil, fmt.Errorf("empty symbol found during extraction")
		}

		symbols = append(symbols, symbol)
	}

	return symbols, nil
}
