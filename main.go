package main

import (
	"embed"
	"github.com/alexflint/go-arg"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/buffers"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/loggers"
	"github.com/phoobynet/trade-ripper/queries"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/pkg/profile"
	"github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var (
	//go:embed dist
	dist              embed.FS
	quitChannel       = make(chan os.Signal, 1)
	options           configuration.Options
	sipReader         *alpaca.TradeReader
	rawMessageChannel = make(chan []byte, 1_000)
	errorsChannel     = make(chan error, 100)
	errorsReceived    = 0
)

func main() {
	defer profile.Start(profile.MemProfileHeap).Stop()

	defer func() {
		loggers.Close()
	}()

	arg.MustParse(&options)

	if options.DBInfluxPort == 0 {
		options.DBInfluxPort = 9009
	}

	if options.DBPostgresPort == 0 {
		options.DBPostgresPort = 8812
	}

	if options.WebServerPort == 0 {
		options.WebServerPort = 3000
	}

	signal.Notify(quitChannel, os.Interrupt)

	questDBErr := queries.InitQuestDB(options)
	if questDBErr != nil {
		panic(questDBErr)
	}

	go run(options)

	server.Run(options, dist)
}

func run(options configuration.Options) {
	logrus.Info("Starting up Trade Reader...")

	questTradeBuffer := buffers.NewQuestBuffer(options)

	sipReader = alpaca.NewTradeReader(&alpaca.TradeReaderConfig{
		Key:               os.Getenv("APCA_API_KEY_ID"),
		Secret:            os.Getenv("APCA_API_SECRET_KEY"),
		Symbols:           []string{"*"},
		RawMessageChannel: rawMessageChannel,
		ErrorsChannel:     errorsChannel,
		Options:           options,
	})

	go func() {
		questTradeBuffer.Start()
	}()

	go func() {
		tradeReaderStartErr := sipReader.Start()

		if tradeReaderStartErr != nil {
			logrus.Fatalln(tradeReaderStartErr)
		}
	}()

	logrus.Info("Trade Reader has started and is waiting for trades...")

	for {
		select {
		case <-quitChannel:
			logrus.Info("Shutting down...")
			_ = sipReader.Stop()
			os.Exit(1)
		case rawMessage := <-rawMessageChannel:
			questTradeBuffer.Add(rawMessage)
		case err := <-errorsChannel:
			errorsReceived++
			logrus.Error(err)
		default:
			logrus.Warn("Overloaded!")
		}
	}
}
