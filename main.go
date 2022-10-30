package main

import (
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

var quitChannel = make(chan os.Signal, 1)
var options configuration.Options
var sipReader *alpaca.TradeReader

var rawMessageChannel = make(chan []byte, 50_000)
var tradesChannel = make(chan alpaca.TradeRow, 100_000)
var errorsChannel = make(chan error, 100)
var errorsReceived = 0
var tradeCount int64

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

	signal.Notify(quitChannel, os.Interrupt)

	questDBErr := queries.InitQuestDB(options)
	if questDBErr != nil {
		panic(questDBErr)
	}

	go run(options)

	server.Run(options)
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
