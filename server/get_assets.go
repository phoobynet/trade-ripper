package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/internal/market/asset"
	"net/http"
	"strings"
)

// getAssets returns a map of assets de-fluffed.
func getAssets(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tickers := r.URL.Query().Get("tickers")

	simplifiedAssets := asset.GetRepositoryInstance().ManySimplified(strings.Split(tickers, ","))

	writeJSONErr := writeJSON(w, http.StatusOK, simplifiedAssets)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
