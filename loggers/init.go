package loggers

import (
	"fmt"
	"github.com/phoobynet/trade-ripper/loggers/hooks"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"io"
	"os"
	"time"
)

var errorLogFile *os.File
var logFile *os.File

func init() {
	now := time.Now().Format("20060102_150405")
	elf, errorLogErr := os.Open(fmt.Sprintf("sipper-ripper_%s_errors.log", now))

	if errorLogErr != nil {
		panic("unable to create error log file")
	}

	errorLogFile = elf

	logrus.AddHook(&writer.Hook{
		Writer: io.MultiWriter(os.Stderr, errorLogFile),
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	lf, logFileErr := os.Open(fmt.Sprintf("sipper-ripper_%s.log", now))

	if logFileErr != nil {
		panic("unable to create log file")
	}

	logFile = lf

	logrus.AddHook(&writer.Hook{
		Writer: io.MultiWriter(os.Stdout, logFile),
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	logrus.AddHook(hooks.NewBroadcastHook())
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
