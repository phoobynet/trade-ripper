package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s *Server) tradeSymbolsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tickers, tickersErr := s.latestTradeRepository.GetKeys()

	if tickersErr != nil {
		_ = WriteErr(w, http.StatusInternalServerError, tickersErr)
		return
	}

	err := writeJSON(w, http.StatusOK, tickers)

	if err != nil {
		_ = WriteErr(w, http.StatusInternalServerError, err)
	}
}
