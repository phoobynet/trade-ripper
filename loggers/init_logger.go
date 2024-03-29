package loggers

import (
	"fmt"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"io"
	"os"
	"time"
)

var (
	errorLogFile *os.File
	logFile      *os.File
)

func InitLogger(webServer *server.Server) {
	now := time.Now().Format("20060102_150405")

	lf, logFileErr := os.Create(fmt.Sprintf("trade_ripper_%s.log", now))

	if logFileErr != nil {
		panic(logFileErr)
	}

	logFile = lf

	// Really important
	logrus.SetOutput(io.Discard)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.AddHook(&writer.Hook{
		Writer: io.MultiWriter(os.Stdout, logFile),
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
		},
	})

	elf, errorLogErr := os.Create(fmt.Sprintf("trade_ripper_%s_errors.log", now))

	if errorLogErr != nil {
		panic(errorLogErr)
	}

	errorLogFile = elf

	logrus.AddHook(&writer.Hook{
		Writer: io.MultiWriter(os.Stderr, errorLogFile),
		LogLevels: []logrus.Level{
			logrus.ErrorLevel,
		},
	})
}

func Close() {
	if errorLogFile != nil {
		_ = errorLogFile.Sync()
		_ = errorLogFile.Close()
	}

	if logFile != nil {
		_ = logFile.Sync()
		_ = logFile.Close()
	}
}
