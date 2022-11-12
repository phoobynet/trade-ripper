package server

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) Listen() {
	listenErr := http.ListenAndServe(fmt.Sprintf(":%d", s.options.WebServerPort), s.router)

	if listenErr != nil {
		log.Fatalln(listenErr)
	}
}
