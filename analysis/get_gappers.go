package analysis

import (
	number_diff "github.com/phoobynet/number-diff"
	"github.com/phoobynet/trade-ripper/internal/market/calendars"
	"github.com/phoobynet/trade-ripper/internal/market/prices"
	"log"
	"sort"
	"time"
)

type Gapper struct {
	Ticker        string  `json:"ticker"`
	PreviousClose float64 `json:"pc"`
	PrevDate      string  `json:"pd"`
	Price         float64 `json:"p"`
	Change        float64 `json:"c"`
	Percent       float64 `json:"cp"`
}

var previousClosingPrices = make(map[string]float64)
var previousCalendar *calendars.Calendar

func GetGappers(latestPrices map[string]float64) []Gapper {
	if previousCalendar == nil {
		c, p := prices.GetPreviousClosingPrices()
		previousCalendar = &p
		previousClosingPrices = c
	}

	now := time.Now().Format("2006-01-02")

	if now != previousCalendar.Date {
		c, p := prices.GetPreviousClosingPrices()
		previousCalendar = &p
		previousClosingPrices = c
	}

	var results []Gapper

	for ticker, price := range latestPrices {
		previousClosingPrice := previousClosingPrices[ticker]

		change, err := number_diff.DiffWithLocale(previousClosingPrice, price, "USD")

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, Gapper{
			Ticker:        ticker,
			PreviousClose: previousClosingPrice,
			PrevDate:      previousCalendar.Date,
			Price:         price,
			Change:        change.Diff,
			Percent:       change.PctDiff,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Percent > results[j].Percent
	})

	return results
}
