package main

import (
	"embed"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/dgraph-io/badger/v3"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/database"
	"github.com/phoobynet/trade-ripper/loggers"
	"github.com/phoobynet/trade-ripper/queries"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/phoobynet/trade-ripper/trades"
	"github.com/phoobynet/trade-ripper/trades/adapters"
	"github.com/phoobynet/trade-ripper/trades/writers"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	//go:embed dist
	dist              embed.FS
	quitChannel       = make(chan os.Signal, 1)
	options           configuration.Options
	sipReader         *alpaca.TradeReader
	rawMessageChannel = make(chan []byte, 1_000_000)
	errorsChannel     = make(chan error, 1)
	errorsReceived    = 0
	tradesChannel     = make(chan []trades.Trade, 10_000)
	tradesBuffer      = make([]trades.Trade, 0)
	tradesWriterLock  = sync.Mutex{}
	latestTradesDB    *badger.DB
	webServer         *server.WebServer
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

	questDBErr := database.StartPostgresConnection(options)
	if questDBErr != nil {
		panic(questDBErr)
	}

	var badgerErr error

	latestTradesDB, badgerErr = badger.Open(badger.DefaultOptions("/latest_trades"))

	if badgerErr != nil {
		panic(badgerErr)
	}

	defer func(latestTradesDB *badger.DB) {
		_ = latestTradesDB.Close()
	}(latestTradesDB)

	webServer = server.NewWebServer(options, dist, latestTradesDB)

	loggers.InitLogger(webServer)

	go run(options)

	webServer.Listen()
}

func run(options configuration.Options) {
	logrus.Info("Starting up Trade Reader...")

	// invoke when we have accumulated enough trades to write to the database
	var tradeWriter writers.TradeWriter
	if options.Class == "crypto" {
		tradeWriter = writers.NewCryptoWriter(options)
	} else {
		tradeWriter = writers.NewUSEquityWriter(options)
	}

	sipReader = alpaca.NewTradeReader(&alpaca.TradeReaderConfig{
		Key:               os.Getenv("APCA_API_KEY_ID"),
		Secret:            os.Getenv("APCA_API_SECRET_KEY"),
		Symbols:           []string{"*"},
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
				tradesWriterLock.Lock()
				defer tradesWriterLock.Unlock()

				l := int64(len(tradesBuffer))

				if l > 0 {
					tradeWriter.Write(tradesBuffer)
					tradesBuffer = make([]trades.Trade, 0)
					count += l
					webServer.Publish(map[string]any{
						"message": "count",
						"data": map[string]any{
							"n": count,
						},
					})
				} else {
					webServer.Publish(map[string]any{
						"message": "ping",
					})
				}
			}()
		}
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
			tradeMessages, adapterErr := adapters.AdaptRawMessageToTrades(rawMessage)

			if adapterErr != nil {
				logrus.Panicf("Error converting raw message to trade: %s", adapterErr)
			}
			tradesChannel <- tradeMessages
		case tradeMessages := <-tradesChannel:
			tradesWriterLock.Lock()
			_ = latestTradesDB.Update(func(txn *badger.Txn) error {
				for _, trade := range tradeMessages {
					_ = txn.Set([]byte(trade["S"].(string)), []byte(fmt.Sprintf("%6.4f,%6.4f,%d", trade["s"].(float64), trade["p"].(float64), trade["t"].(int64))))
				}
				return nil
			})
			tradesBuffer = append(tradesBuffer, tradeMessages...)
			tradesWriterLock.Unlock()
		case err := <-errorsChannel:
			errorsReceived++
			logrus.Error(err)
		}
	}
}
