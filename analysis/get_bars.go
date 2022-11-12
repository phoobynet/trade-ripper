package analysis

import (
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"strings"
	"time"
)

type Bar struct {
	Open      float64   `json:"o"`
	High      float64   `json:"h"`
	Low       float64   `json:"l"`
	Close     float64   `json:"c"`
	Volume    float64   `json:"v"`
	Timestamp time.Time `json:"t"`
}

const getBarsSQL = `
select 
	first(price) open, 
	max(price) high, 
	min(price) low, 
	last(price) close, 
	sum(size) volume, 
	last(timestamp) timestamp
from us_equity
where 
	ticker = ? and timestamp in ':date'
sample by :interval FILL(prev)
align to CALENDAR with offset '00:00'
`

func GetBars(ticker string, date string, interval string) ([]Bar, error) {
	db, dbErr := postgres.Get()
	if dbErr != nil {
		return nil, dbErr
	}

	var bars []Bar

	finalSQL := strings.Replace(getBarsSQL, ":date", date, 1)
	finalSQL = strings.Replace(finalSQL, ":interval", interval, 1)

	result := db.Raw(sql, ticker)
	scanErr := result.Scan(&bars).Error

	if scanErr != nil {
		return nil, scanErr
	}

	return bars, nil
}
