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
	router                *httprouter.Router
}

func NewServer(options configuration.Options, dist embed.FS, latestTradeRepository *tradeskv.LatestTradeRepository) *Server {
	sseServer := sse.New()
	sseServer.CreateStream("events")

	router := httprouter.New()
	webServer := &Server{
		options,
		latestTradeRepository,
		sseServer,
		router,
	}

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.GET("/events", webServer.eventsHandler)
	router.GET("/trades/latest", webServer.tradesLatestHandler)
	router.GET("/trades/symbols", webServer.tradeSymbolsHandler)
	router.GET("/classHandler", webServer.classHandler)
	router.GET("/bars/:ticker/:date", webServer.barsHandler)

	fsys, distFSErr := fs.Sub(dist, "dist")

	if distFSErr != nil {
		panic(distFSErr)
	}

	router.ServeFiles("/", http.FS(fsys))

	return webServer
}
