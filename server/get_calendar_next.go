package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/internal/market/calendars"
	"net/http"
)

func getCalendarNext(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	next := calendars.Next()

	writeJSONErr := writeJSON(w, http.StatusOK, next)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
