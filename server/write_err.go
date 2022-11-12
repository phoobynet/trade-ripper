package server

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func WriteErr(w http.ResponseWriter, statusCode int, err error) error {
	logrus.Error(err)
	return writeJSON(w, statusCode, map[string]string{
		"error": err.Error(),
	})
}
