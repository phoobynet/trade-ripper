package server

import (
	"encoding/json"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/r3labs/sse/v2"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var server *sse.Server

func Publish(message any) {
	if server == nil {
		logrus.Error("web server is nil")
		return
	}

	if message == nil {
		logrus.Error("message is nil")
		return
	}

	data, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		return
	}

	server.Publish("messages", &sse.Event{
		Data: data,
	})
}

func Run(options configuration.Options) {
	server = sse.New()
	server.CreateStream("messages")

	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			println("The client is disconnected here")
			return
		}()

		fmt.Printf("Client connected...\n")

		server.ServeHTTP(w, r)
	})

	mux.HandleFunc("/api/class", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(fmt.Sprintf(`{"class": "%s"}`, options.Class)))
	})

	staticFilesServer := http.FileServer(http.Dir("./public"))

	mux.Handle("/", staticFilesServer)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", options.WebServerPort), c.Handler(mux))

	if err != nil {
		log.Fatalln(err)
	}
}
