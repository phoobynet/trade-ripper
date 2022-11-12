package server

import (
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/tradeskv"
	"github.com/r3labs/sse/v2"
	"io/fs"
	"net/http"
)

type Server struct {
	options               configuration.Options
	latestTradeRepository *tradeskv.LatestTradeRepository
	sseServer             *sse.Server
	mux                   *http.ServeMux
}

func NewServer(options configuration.Options, dist embed.FS, latestTradeRepository *tradeskv.LatestTradeRepository) *Server {
	sseServer := sse.New()
	sseServer.CreateStream("events")

	webServer := &Server{
		options,
		latestTradeRepository,
		sseServer,
		http.NewServeMux(),
	}

	fsys, distFSErr := fs.Sub(dist, "dist")

	if distFSErr != nil {
		panic(distFSErr)
	}

	router := httprouter.New()
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		w.WriteHeader(http.StatusNoContent)
	})
	router.GET("/api/events", webServer.eventsHandler)
	router.GET("/api/trades/latest", webServer.tradesLatestHandler)
	router.GET("/api/trades/symbols", webServer.tradeSymbolsHandler)
	router.GET("/api/bars/:ticker/:interval/:date", webServer.barsHandler)
	router.GET("/api/class", webServer.classHandler)

	webServer.mux.Handle("/", http.FileServer(http.FS(fsys)))
	webServer.mux.Handle("/api/", router)

	return webServer
}
