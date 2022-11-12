package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s *Server) classHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := writeJSON(w, http.StatusOK, map[string]any{
		"classHandler": s.options.Class,
	})

	if err != nil {
		_ = WriteErr(w, http.StatusInternalServerError, err)
		return
	}
}
