package server

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	j, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(j)

	return err
}
