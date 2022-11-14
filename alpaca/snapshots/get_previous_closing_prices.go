package snapshots

import (
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/sirupsen/logrus"
)

func GetPreviousClosingPrices() (map[string]float64, calendars.Calendar) {
	db := localdb.Get()

	previousCalendar := calendars.Previous()

	var previousCloses []PreviousClose

	err := db.Raw("select ticker, price, date from previous_close where date = ?", previousCalendar.Date).Scan(&previousCloses).Error

	if err != nil {
		logrus.Fatal(err)
	}

	previousClosingPrices := make(map[string]float64, 0)

	for _, previousClose := range previousCloses {
		previousClosingPrices[previousClose.Ticker] = previousClose.Price
	}

	return previousClosingPrices, *previousCalendar
}
