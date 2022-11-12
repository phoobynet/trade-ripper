package server

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s *Server) eventsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	go func() {
		// Received Browser Disconnection
		<-r.Context().Done()
		println("The client is disconnected here")
	}()

	fmt.Printf("Client connected...\n")

	s.sseServer.ServeHTTP(w, r)
}
