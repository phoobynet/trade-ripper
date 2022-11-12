package server

import (
	_ "embed"
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/analysis"
	"net/http"
)

func getBars(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ticker := p.ByName("ticker")
	date := p.ByName("date")
	interval := p.ByName("interval")

	bars, barsErr := analysis.GetBars(ticker, date, interval)

	if barsErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, barsErr)
		return
	}

	writeJSONErr := writeJSON(w, http.StatusOK, bars)
	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
