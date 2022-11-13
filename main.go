package main

import (
	"embed"
	"github.com/alexflint/go-arg"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/alpaca/assets"
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/loggers"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/phoobynet/trade-ripper/tradesdb"
	"github.com/phoobynet/trade-ripper/tradesdb/adapters"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"github.com/phoobynet/trade-ripper/tradesdb/queries"
	"github.com/phoobynet/trade-ripper/tradesdb/writers"
	"github.com/phoobynet/trade-ripper/tradeskv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	//go:embed dist
	dist                  embed.FS
	quitChannel           = make(chan os.Signal, 1)
	options               configuration.Options
	tradeReader           *alpaca.TradeReader
	rawMessageChannel     = make(chan []byte, 1_000_000)
	errorsChannel         = make(chan error, 1)
	errorsReceived        = 0
	tradesChannel         = make(chan []tradesdb.Trade, 10_000)
	tradesBuffer          = make([]tradesdb.Trade, 0)
	tradesWriterLock      = sync.Mutex{}
	latestTradeRepository *tradeskv.LatestTradeRepository
	webServer             *server.Server
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

	postgres.Initialize(options)
	assets.Initialize()
	calendars.Initialize()

	latestTradeRepository = tradeskv.NewLatestRepository(options)

	defer func(repository *tradeskv.LatestTradeRepository) {
		repository.Close()
	}(latestTradeRepository)

	webServer = server.NewServer(options, dist, latestTradeRepository)

	loggers.InitLogger(webServer)

	logrus.Info("Starting up Trade Reader...")

	// invoke when we have accumulated enough trades to write to the database
	tradeWriter := writers.CreateTradeWriter(options)

	tickers := options.ExtractTickers()

	tradeReader = alpaca.NewTradeReader(&alpaca.TradeReaderConfig{
		Key:               os.Getenv("APCA_API_KEY_ID"),
		Secret:            os.Getenv("APCA_API_SECRET_KEY"),
		Symbols:           tickers,
		RawMessageChannel: rawMessageChannel,
		ErrorsChannel:     errorsChannel,
		Options:           options,
	})

	go func() {
		count, countErr := queries.Count(options)

		if countErr != nil {
			logrus.Panicf("Error counting trades: %s", countErr)
		}
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logrus.Errorf("Recovered in f %v", r)
					}
					time.Sleep(1 * time.Second)
				}()
				tradesWriterLock.Lock()
				defer tradesWriterLock.Unlock()

				l := int64(len(tradesBuffer))

				if l > 0 {
					tradeWriter.Write(tradesBuffer)
					tradesBuffer = make([]tradesdb.Trade, 0)
					count += l
					server.PublishEvent(map[string]any{
						"message": "count",
						"data": map[string]any{
							"n": count,
						},
					})
				} else {
					server.PublishEvent(map[string]any{
						"message": "ping",
					})
				}
			}()
		}
	}()

	go func() {
		logrus.Info("Trade Reader has started and is waiting for trades...")
		tradeReaderStartErr := tradeReader.Start()

		if tradeReaderStartErr != nil {
			logrus.Fatalln(tradeReaderStartErr)
		}
	}()

	go func() {
		for {
			select {
			case <-quitChannel:
				logrus.Info("Shutting down...")
				_ = tradeReader.Stop()
				os.Exit(1)
			case rawMessage := <-rawMessageChannel:
				tradeMessages, adapterErr := adapters.AdaptRawMessageToTrades(rawMessage)

				if adapterErr != nil {
					logrus.Panicf("Error converting raw message to trade: %s", adapterErr)
				}
				tradesChannel <- tradeMessages
			case tradeMessages := <-tradesChannel:
				tradesWriterLock.Lock()
				tradesBuffer = append(tradesBuffer, tradeMessages...)
				latestTradeUpdateErr := latestTradeRepository.Update(tradeMessages)
				if latestTradeUpdateErr != nil {
					logrus.Error(latestTradeUpdateErr)
				}
				tradesWriterLock.Unlock()
			case err := <-errorsChannel:
				errorsReceived++
				logrus.Error(err)
			}
		}
	}()

	webServer.Listen()
}
