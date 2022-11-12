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

func (s *Server) barsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ticker := p.ByName("ticker")
	date := p.ByName("date")
	db, dbErr := postgres.Get()

	if dbErr != nil {
		_ = WriteErr(w, http.StatusInternalServerError, dbErr)
		return
	}

	var bars []Bar

	// HACK: Had to find/replace "timestamp in " clause
	sql := strings.Replace(`
		select 
			first(price) open, 
			max(price) high, 
			min(price) low, 
			latest(price) close, 
			sum(size) volume, 
			latest(timestamp) timestamp
		where 
			ticker = ? and timestamp in ':today'
	`, ":today", date, 1)

	result := db.Raw(sql, ticker)
	scanErr := result.Scan(&bars).Error

	if scanErr != nil {
		_ = WriteErr(w, http.StatusInternalServerError, scanErr)
		return
	}

	writeJSONErr := writeJSON(w, http.StatusOK, bars)
	if writeJSONErr != nil {
		_ = WriteErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
