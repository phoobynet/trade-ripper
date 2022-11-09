package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/tradeskv"
	"github.com/r3labs/sse/v2"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

type WebServer struct {
	options               configuration.Options
	latestTradeRepository *tradeskv.LatestTradeRepository
	sseServer             *sse.Server
	mux                   *http.ServeMux
	corsOptions           *cors.Cors
}

func NewWebServer(options configuration.Options, dist embed.FS, latestTradeRepository *tradeskv.LatestTradeRepository) *WebServer {
	sseServer := sse.New()
	sseServer.CreateStream("events")

	mux := http.NewServeMux()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	webServer := &WebServer{
		options,
		latestTradeRepository,
		sseServer,
		mux,
		corsOptions,
	}

	mux.HandleFunc("/events", webServer.eventsHandler)
	mux.HandleFunc("/trades/latest", webServer.tradesLatest)
	mux.HandleFunc("/trades/symbols", webServer.tradesSymbols)
	mux.HandleFunc("/api/class", webServer.class)

	fsys, distFSErr := fs.Sub(dist, "dist")

	if distFSErr != nil {
		panic(distFSErr)
	}

	staticFilesServer := http.FileServer(http.FS(fsys))

	mux.Handle("/", staticFilesServer)

	return webServer
}

func (ws *WebServer) PublishEvent(message any) {
	if ws.sseServer == nil {
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

	ws.sseServer.Publish("events", &sse.Event{
		Data: data,
	})
}

func (ws *WebServer) eventsHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		// Received Browser Disconnection
		<-r.Context().Done()
		println("The client is disconnected here")
	}()

	fmt.Printf("Client connected...\n")

	ws.sseServer.ServeHTTP(w, r)
}

func (ws *WebServer) tradesLatest(w http.ResponseWriter, r *http.Request) {
	tickersQuery := r.URL.Query().Get("tickers")

	if tickersQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		tickers := strings.Split(tickersQuery, ",")
		trades, latestTradeErr := ws.latestTradeRepository.Get(tickers)

		if latestTradeErr != nil {
			logrus.Error(latestTradeErr)
			_, _ = w.Write([]byte(latestTradeErr.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		j, jErr := json.Marshal(trades)
		if jErr != nil {
			logrus.Error(jErr)
			_, _ = w.Write([]byte(jErr.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, writeErr := w.Write(j)

		if writeErr != nil {
			logrus.Error(writeErr)
			_, _ = w.Write([]byte(writeErr.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (ws *WebServer) class(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(fmt.Sprintf(`{"class": "%s"}`, ws.options.Class)))
}

func (ws *WebServer) Listen() {
	listenErr := http.ListenAndServe(fmt.Sprintf(":%d", ws.options.WebServerPort), ws.corsOptions.Handler(ws.mux))

	if listenErr != nil {
		log.Fatalln(listenErr)
	}
}

func (ws *WebServer) tradesSymbols(w http.ResponseWriter, r *http.Request) {
	tickers, tickersErr := ws.latestTradeRepository.GetKeys()

	if tickersErr != nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tickersErr.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	j, jErr := json.Marshal(tickers)

	if jErr != nil {
		logrus.Error(jErr)
	}

	_, writeErr := w.Write(j)
	if writeErr != nil {
		logrus.Error(writeErr)
	}
}
