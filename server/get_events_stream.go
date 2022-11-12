package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var sseServer *sse.Server
var sseMutex sync.Mutex

func initSSE() {
	sseMutex.Lock()
	defer sseMutex.Unlock()

	if sseServer == nil {
		sseServer = sse.New()
		sseServer.CreateStream("events")
	}
}
func init() {
	initSSE()
}

func getEventsStream(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	go func() {
		// Received Browser Disconnection
		<-r.Context().Done()
		println("The client is disconnected here")
	}()

	fmt.Printf("Client connected...\n")

	sseServer.ServeHTTP(w, r)
}

func PublishEvent(message any) {
	if sseServer == nil {
		logrus.Warn("SSE Server is not initialized")
		return
	}

	if message == nil {
		logrus.Warn("Cannot send nil message")
		return
	}

	data, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		return
	}

	sseServer.Publish("events", &sse.Event{
		Data: data,
	})
}
