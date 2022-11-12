package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"net/http"
	"strings"
)

type VolumeLeader struct {
	Ticker string `json:"ticker"`
	Volume string `json:"volume"`
	Price  string `json:"price"`
}

func getVolumeLeaders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db, dbErr := postgres.Get()

	if dbErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, dbErr)
		return
	}

	limit := ps.ByName("limit")
	date := ps.ByName("date")

	const sql = `
		select ticker, sum(size) volume, last (price) price
    from
        us_equity
    where
        timestamp in ':date'
    group by
        ticker
    order by
        volume desc
        limit :limit
	`

	finalSQL := strings.Replace(sql, ":date", date, 1)
	finalSQL = strings.Replace(finalSQL, ":limit", limit, 1)

	var volumeLeaders []VolumeLeader

	scanErr := db.Raw(finalSQL).Scan(&volumeLeaders).Error

	if scanErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, scanErr)
		return
	}

	writeJSONErr := writeJSON(w, http.StatusOK, volumeLeaders)

	if writeJSONErr != nil {
		_ = writeErr(w, http.StatusInternalServerError, writeJSONErr)
	}
}
