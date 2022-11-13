package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/r3labs/sse/v2"
	"net/http"
	"sync"
)

var sseServer *sse.Server
var sseMutex sync.Mutex

func InitSSE() {
	sseMutex.Lock()
	defer sseMutex.Unlock()

	if sseServer == nil {
		sseServer = sse.New()
		sseServer.CreateStream("events")
	}
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
	// DON'T DO THIS: logrus.* functions will cause infinite recursion
	if sseServer == nil {
		panic("SSE Server is not initialized")
	}

	if message == nil {
		panic("Cannot send nil message")
	}

	data, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		return
	}

	sseServer.Publish("events", &sse.Event{
		Data: data,
	})
}
