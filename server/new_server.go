package server

import (
	"embed"
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/trade-ripper/internal/configuration"
	"github.com/rs/cors"
	"io/fs"
	"net/http"
)

type Server struct {
	options configuration.Options
	mux     *http.ServeMux
}

func NewServer(options configuration.Options, dist embed.FS) *Server {
	webServer := &Server{
		options,
		http.NewServeMux(),
	}

	fsys, distFSErr := fs.Sub(dist, "dist")

	if distFSErr != nil {
		panic(distFSErr)
	}

	router := httprouter.New()

	myCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	router.GET("/api/events", getEventsStream)
	router.GET("/api/class", createGetClassHandler(options.Class))

	router.GET("/api/assets", getAssets)
	router.GET("/api/bars/:ticker/:date/:interval", getBars)
	router.GET("/api/volume-leaders/:date/:limit", getVolumeLeaders)
	router.GET("/api/calendar/previous", getCalendarPrevious)
	router.GET("/api/calendar/current", getCalendarCurrent)
	router.GET("/api/calendar/next", getCalendarNext)
	router.GET("/api/market-status", getMarketStatus)

	webServer.mux.Handle("/", http.FileServer(http.FS(fsys)))
	webServer.mux.Handle("/api/", myCors.Handler(router))

	return webServer
}
