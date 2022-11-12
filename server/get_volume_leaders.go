package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/analysis"
	"net/http"
	"strconv"
)

func getVolumeLeaders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limitParam := ps.ByName("limit")
	date := ps.ByName("date")

	limit, limitErr := strconv.Atoi(limitParam)

	if limitErr != nil {
		_ = writeErr(w, http.StatusBadRequest, limitErr)
		return
	}

	volumeLeaders, volumeLeadersErr := analysis.GetVolumeLeaders(date, limit)

	if volumeLeadersErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, volumeLeadersErr)
		return
	}

	writeJSONErr := writeJSON(w, http.StatusOK, volumeLeaders)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
