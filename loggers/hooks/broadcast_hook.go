package hooks

import (
	"fmt"
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
	_, err := entry.Bytes()
	if err != nil {
		return err
	}

	logMessage := make(map[string]any)
	logMessage["type"] = entry.Level.String()
	logMessage["msg"] = entry.Message
	logMessage["tradeTrades"] = entry.Data["totalTrades"]
	logMessage["time"] = entry.Time
	fmt.Printf("%+v", logMessage)

	server.Broadcast(logMessage)

	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *BroadcastHook) Levels() []logrus.Level {
	return hook.LogLevels
}
