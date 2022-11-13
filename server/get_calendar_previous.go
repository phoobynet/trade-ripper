package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"net/http"
)

func getCalendarPrevious(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	previous := calendars.Previous()

	writeJSONErr := writeJSON(w, http.StatusOK, previous)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
