package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/analysis"
	"net/http"
)

func createGetGappersHandler(latestPrices map[string]float64) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		gappers := analysis.GetGappers(latestPrices)
		writeJSONErr := writeJSON(w, http.StatusOK, gappers)

		if writeJSONErr != nil {
			_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
		}
	}
}
