package hooks

import (
	"github.com/phoobynet/trade-ripper/server"
	"github.com/sirupsen/logrus"
	"io"
)

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
	server.Broadcast(server.Message{
		MessageType:    entry.Level.String(),
		MessageContent: string(line),
	})

	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *BroadcastHook) Levels() []logrus.Level {
	return hook.LogLevels
}
