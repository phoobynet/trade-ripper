package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/r3labs/sse/v2"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"io/fs"
	"log"
	"net/http"
)

var (
	server *sse.Server
)

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

func Run(options configuration.Options, dist embed.FS) {
	server = sse.New()
	server.CreateStream("messages")

	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			println("The client is disconnected here")
		}()

		fmt.Printf("Client connected...\n")

		server.ServeHTTP(w, r)
	})

	mux.HandleFunc("/api/class", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(fmt.Sprintf(`{"class": "%s"}`, options.Class)))
	})

	fmt.Println("using embed mode")
	fsys, distFSErr := fs.Sub(dist, "dist")

	if distFSErr != nil {
		panic(distFSErr)
	}

	staticFilesServer := http.FileServer(http.FS(fsys))

	mux.Handle("/", staticFilesServer)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	listenErr := http.ListenAndServe(fmt.Sprintf(":%d", options.WebServerPort), c.Handler(mux))

	if listenErr != nil {
		log.Fatalln(listenErr)
	}
}
