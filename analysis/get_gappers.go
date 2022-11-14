package analysis

import (
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"github.com/phoobynet/trade-ripper/alpaca/snapshots"
	"github.com/phoobynet/trade-ripper/utils"
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
		c, p := snapshots.GetPreviousClosingPrices()
		previousCalendar = &p
		previousClosingPrices = c
	}

	now := time.Now().Format("2006-01-02")

	if now != previousCalendar.Date {
		c, p := snapshots.GetPreviousClosingPrices()
		previousCalendar = &p
		previousClosingPrices = c
	}

	var results []Gapper

	if latestPrices == nil || len(latestPrices) == 0 {
		panic("latestPrices is nil or empty")
	}

	for ticker, price := range latestPrices {
		previousClosingPrice := previousClosingPrices[ticker]

		change := utils.NumberDiff(previousClosingPrice, price)

		results = append(results, Gapper{
			Ticker:        ticker,
			PreviousClose: previousClosingPrice,
			PrevDate:      previousCalendar.Date,
			Price:         price,
			Change:        change.CashDifference,
			Percent:       change.PercentDifference,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Percent > results[j].Percent
	})

	return results
}
