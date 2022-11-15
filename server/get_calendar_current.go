package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/internal/market/calendars"
	"net/http"
)

func getCalendarCurrent(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	current := calendars.Current()

	if current == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		writeJSONErr := writeJSON(w, http.StatusOK, current)

		if writeJSONErr != nil {
			_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
		}
	}
}
