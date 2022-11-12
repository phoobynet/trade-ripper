package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func createGetClassHandler(class string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := writeJSON(w, http.StatusOK, map[string]any{
			"class": class,
		})

		if err != nil {
			_ = writeErr(w, http.StatusInternalServerError, err)
			return
		}
	}
}
