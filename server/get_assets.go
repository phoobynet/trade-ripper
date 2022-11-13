package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/alpaca/assets"
	"net/http"
	"strings"
)

// getAssets returns a map of assets de-fluffed.
func getAssets(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tickers := r.URL.Query().Get("tickers")

	writeJSONErr := writeJSON(w, http.StatusOK, assets.ManySimplified(strings.Split(tickers, ",")))

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
