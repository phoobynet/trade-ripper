package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/internal/market/status"
	"net/http"
)

func getMarketStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	marketStatus := status.GetMarketStatus()

	writeJSONErr := writeJSON(w, http.StatusOK, marketStatus)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
