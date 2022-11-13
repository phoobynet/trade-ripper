package server

import (
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/tradeskv"
	"io/fs"
	"net/http"
)

type Server struct {
	options               configuration.Options
	latestTradeRepository *tradeskv.LatestTradeRepository
	mux                   *http.ServeMux
}

func NewServer(options configuration.Options, dist embed.FS, latestTradeRepository *tradeskv.LatestTradeRepository) *Server {
	webServer := &Server{
		options,
		latestTradeRepository,
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
	router.GET("/api/events", getEventsStream)
	router.GET("/api/class", createGetClassHandler(options.Class))

	router.GET("/api/assets", getAssets)
	router.GET("/api/bars/:ticker/:date/:interval", getBars)
	router.GET("/api/volume-leaders/:date/:limit", getVolumeLeaders)
	router.GET("/api/calendar/previous", getCalendarPrevious)
	router.GET("/api/calendar/current", getCalendarCurrent)

	webServer.mux.Handle("/", http.FileServer(http.FS(fsys)))
	webServer.mux.Handle("/api/", router)

	return webServer
}
