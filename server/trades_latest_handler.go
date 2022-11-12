package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (s *Server) tradesLatestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tickersQuery := r.URL.Query().Get("tickers")

	if tickersQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		trades, latestTradeErr := s.latestTradeRepository.Get(strings.Split(tickersQuery, ","))

		if latestTradeErr != nil {
			_ = WriteErr(w, http.StatusInternalServerError, latestTradeErr)
			return
		}

		writeJSONErr := writeJSON(w, http.StatusOK, trades)

		if writeJSONErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
