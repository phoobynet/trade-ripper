package server

import (
	_ "embed"
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"net/http"
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

func getBars(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ticker := p.ByName("ticker")
	date := p.ByName("date")
	interval := p.ByName("interval")
	db, dbErr := postgres.Get()

	if dbErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, dbErr)
		return
	}

	var bars []Bar

	// HACK: Had to find/replace "timestamp in " clause
	sql := strings.ReplaceAll(getBarsSQL, ":date", date)
	sql = strings.ReplaceAll(sql, ":interval", interval)

	result := db.Raw(sql, ticker)
	scanErr := result.Scan(&bars).Error

	if scanErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, scanErr)
		return
	}

	writeJSONErr := writeJSON(w, http.StatusOK, bars)
	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
