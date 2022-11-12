package server

import (
	"encoding/json"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

func (s *Server) PublishEvent(message any) {
	if s.sseServer == nil {
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

	s.sseServer.Publish("events", &sse.Event{
		Data: data,
	})
}
