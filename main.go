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
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var (
	//go:embed dist
	dist              embed.FS
	quitChannel       = make(chan os.Signal, 1)
	options           configuration.Options
	sipReader         *alpaca.TradeReader
	rawMessageChannel = make(chan []byte, 50_000)
	tradesChannel     = make(chan alpaca.TradeRow, 100_000)
	errorsChannel     = make(chan error, 100)
	errorsReceived    = 0
	tradeCount        int64
)

func main() {
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
			rows, err := alpaca.ConvertToTradeRows(rawMessage)

			if err != nil {
				panic(err)
			}

			for _, row := range rows {
				tradesChannel <- row
			}
		case trade := <-tradesChannel:
			questTradeBuffer.Add(trade)
			tradeCount++
		case err := <-errorsChannel:
			errorsReceived++
			logrus.Error(err)
		}
	}
}
