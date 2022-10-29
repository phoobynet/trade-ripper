package main

import (
	"github.com/alexflint/go-arg"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/buffers"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/loggers"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

var quitChannel = make(chan os.Signal, 1)
var options configuration.Options
var sipReader *alpaca.TradeReader

var rawMessageChannel chan []byte
var tradesChannel chan alpaca.TradeRow
var errorsChannel chan error
var errorsReceived int
var restartsInPeriod int
var restarts int
var lastRestartTime time.Time
var tradeCount int

func main() {
	defer func() {
		loggers.Close()
	}()

	arg.MustParse(&options)
	signal.Notify(quitChannel, os.Interrupt)

	rawMessageChannel = make(chan []byte, 50_000)
	tradesChannel = make(chan alpaca.TradeRow, 100_000)
	errorsChannel = make(chan error, 100)
	errorsReceived = 0
	restartsInPeriod = 0
	restarts = 0

	go run(&options)

	server.Run(&options)
}

func run(options *configuration.Options) {
	logrus.Info("Starting up SIP Reader...")

	questTradeBuffer := buffers.NewQuestBuffer(options)

	sipReader = alpaca.NewTradeReader(&alpaca.TradeReaderConfig{
		Key:               os.Getenv("APCA_API_KEY_ID"),
		Secret:            os.Getenv("APCA_API_SECRET_KEY"),
		Symbols:           []string{"*"},
		RawMessageChannel: rawMessageChannel,
		ErrorsChannel:     errorsChannel,
	})

	go func() {
		questTradeBuffer.Start()
	}()

	go func() {
		sipReaderStartErr := sipReader.Start()

		if sipReaderStartErr != nil {
			logrus.Fatalln(sipReaderStartErr)
		}
	}()

	logrus.Info("SIP Reader has started and is waiting for trades...")

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
