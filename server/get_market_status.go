package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/market"
	"net/http"
)

func getMarketStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	marketStatus := market.GetStatus()

	writeJSONErr := writeJSON(w, http.StatusOK, marketStatus)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
