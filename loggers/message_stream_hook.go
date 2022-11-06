package loggers

import (
	"github.com/sirupsen/logrus"
	"io"
)

import (
	"github.com/phoobynet/trade-ripper/server"
)

type MessageStreamHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	webServer *server.WebServer
}

func NewMessageStreamHook(webServer *server.WebServer) *MessageStreamHook {
	return &MessageStreamHook{
		Writer:    io.Discard,
		LogLevels: logrus.AllLevels,
		webServer: webServer,
	}
}

func (hook *MessageStreamHook) Fire(entry *logrus.Entry) error {
	_, err := entry.Bytes()
	if err != nil {
		return err
	}
	hook.webServer.Publish(map[string]interface{}{
		"type":    entry.Level.String(),
		"message": entry.Message,
		"time":    entry.Time,
		"data":    entry.Data,
	})

	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *MessageStreamHook) Levels() []logrus.Level {
	return hook.LogLevels
}
