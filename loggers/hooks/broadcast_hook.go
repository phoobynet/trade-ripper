package hooks

import (
	"github.com/phoobynet/trade-ripper/server"
	"github.com/sirupsen/logrus"
	"io"
)

var errorCount int

type BroadcastHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func NewBroadcastHook() *BroadcastHook {
	return &BroadcastHook{
		Writer:    io.Discard,
		LogLevels: logrus.AllLevels,
	}
}

func (hook *BroadcastHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	switch entry.Level.String() {
	case "panic":
		fallthrough
	case "fatal":
		fallthrough
	case "error":
		errorCount += 1
		server.Broadcast(server.ErrorMessage{
			Message: server.Message{Type: "error"},
			Msg:     string(line),
			Count:   errorCount,
		})
	default:
		server.Broadcast(server.InfoMessage{
			Message: server.Message{Type: "info"},
			Msg:     string(line),
		})
	}

	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *BroadcastHook) Levels() []logrus.Level {
	return hook.LogLevels
}
