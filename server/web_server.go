package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/r3labs/sse/v2"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WebServer struct {
	options     configuration.Options
	db          *badger.DB
	sseServer   *sse.Server
	mux         *http.ServeMux
	corsOptions *cors.Cors
}

func NewWebServer(options configuration.Options, dist embed.FS, db *badger.DB) *WebServer {
	sseServer := sse.New()
	sseServer.CreateStream("messages")

	mux := http.NewServeMux()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	webServer := &WebServer{
		options,
		db,
		sseServer,
		mux,
		corsOptions,
	}

	mux.HandleFunc("/events", webServer.eventsHandler)
	mux.HandleFunc("/trades/latest", webServer.tradesLatest)
	mux.HandleFunc("/api/class", webServer.tradesLatest)

	fmt.Println("using embed mode")
	fsys, distFSErr := fs.Sub(dist, "dist")

	if distFSErr != nil {
		panic(distFSErr)
	}

	staticFilesServer := http.FileServer(http.FS(fsys))

	mux.Handle("/", staticFilesServer)

	return webServer
}

func (ws *WebServer) Publish(message any) {
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

	ws.sseServer.Publish("messages", &sse.Event{
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
	symbols := r.URL.Query().Get("symbols")

	if symbols == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		ws.db.View(func(txn *badger.Txn) error {
			trades := make(map[string]any)

			for _, symbol := range strings.Split(symbols, ",") {
				trade, err := txn.Get([]byte(strings.ToUpper(symbol)))
				if err != nil {
					if err == badger.ErrKeyNotFound {
						continue
					} else {
						return err
					}
				}

				tradeErr := trade.Value(func(val []byte) error {
					tokens := strings.Split(string(val), ",")
					size, sizeErr := strconv.ParseFloat(tokens[0], 64)
					if sizeErr != nil {
						return sizeErr
					}
					price, priceErr := strconv.ParseFloat(tokens[1], 64)

					if priceErr != nil {
						return priceErr
					}

					timestamp, timestampErr := strconv.ParseInt(tokens[2], 10, 64)

					if timestampErr != nil {
						return timestampErr
					}

					if ws.options.Class == "crypto" {
						trades[symbol] = map[string]any{
							"size":      size,
							"price":     price,
							"timestamp": time.Unix(0, timestamp),
							"tks":       tokens[3],
						}
					} else {
						trades[symbol] = map[string]any{
							"size":      size,
							"price":     price,
							"timestamp": time.Unix(0, timestamp),
						}
					}

					return nil
				})

				if tradeErr != nil {
					logrus.Error(tradeErr)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			j, jErr := json.Marshal(trades)
			if jErr != nil {
				return jErr
			}

			_, writeErr := w.Write(j)

			return writeErr
		})
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
