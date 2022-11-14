package analysis

import (
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"strings"
	"time"
)

const getLatestPricesSQL = `select ticker, round_half_even(last (price), 2) price from us_equity where timestamp in ':date' group by ticker`

func GetLatestPrices(date time.Time) map[string]float64 {
	db, err := postgres.Get()

	if err != nil {
		panic(err)
	}

	finalSQL := strings.Replace(getLatestPricesSQL, ":date", date.Format("2006-01-02"), 1)

	var latestPrices []struct {
		Ticker string
		Price  float64
	}

	queryErr := db.Raw(finalSQL).Scan(&latestPrices).Error

	if queryErr != nil {
		panic(queryErr)
	}

	latestPricesMap := make(map[string]float64)

	for _, latestPrice := range latestPrices {
		latestPricesMap[latestPrice.Ticker] = latestPrice.Price
	}

	return latestPricesMap
}
