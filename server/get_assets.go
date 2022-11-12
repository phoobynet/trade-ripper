package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/alpaca"
	"net/http"
	"strings"
)

func getAssets(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tickers := r.URL.Query().Get("tickers")

	assets, assetErr := alpaca.GetAssets(strings.Split(tickers, ","))

	if assetErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, assetErr)
		return
	}

	writeJSONErr := writeJSON(w, http.StatusOK, assets)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
