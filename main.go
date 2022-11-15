package main

import (
	"embed"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/alpaca/assets"
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"github.com/phoobynet/trade-ripper/alpaca/snapshots"
	"github.com/phoobynet/trade-ripper/analysis"
	"github.com/phoobynet/trade-ripper/internal/configuration"
	_ "github.com/phoobynet/trade-ripper/internal/configuration"
	"github.com/phoobynet/trade-ripper/internal/loggers"
	_ "github.com/phoobynet/trade-ripper/internal/loggers"
	"github.com/phoobynet/trade-ripper/market"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/phoobynet/trade-ripper/tradesdb"
	"github.com/phoobynet/trade-ripper/tradesdb/adapters"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"github.com/phoobynet/trade-ripper/tradesdb/queries"
	"github.com/phoobynet/trade-ripper/tradesdb/writers"
	"github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	//go:embed dist
	dist              embed.FS
	quitChannel       = make(chan os.Signal, 1)
	tradeReader       *alpaca.TradeReader
	rawMessageChannel = make(chan []byte, 100_000)
	errorsChannel     = make(chan error, 1)
	errorsReceived    = 0
	tradesChannel     = make(chan []tradesdb.Trade, 10_000)
	tradesBuffer      = make([]tradesdb.Trade, 0)
	tradesWriterLock  = sync.RWMutex{}
	webServer         *server.Server
	latestPrices      = make(map[string]float64)
)

func main() {
	defer func() {
		loggers.Close()
	}()

	signal.Notify(quitChannel, os.Interrupt)

	tradeRipperOptions := configuration.GetTradeRipperOptions()
	tickers := strings.Split(tradeRipperOptions.Tickers, ",")

	postgres.Initialize(tradeRipperOptions)
	assets.Initialize()
	calendars.Initialize()
	server.InitSSE()

	webServer = server.NewServer(tradeRipperOptions, dist)

	logrus.Info("Starting up Trade Reader...")

	// invoke when we have accumulated enough trades to write to the database
	tradeWriter := writers.CreateTradeWriter(tradeRipperOptions)

	snapshots.CachePreviousClose(tickers)

	for ticker, price := range analysis.GetLatestPrices(time.Now()) {
		latestPrices[ticker] = price
	}

	tradeReader = alpaca.NewTradeReader(&alpaca.TradeReaderConfig{
		Key:               os.Getenv("APCA_API_KEY_ID"),
		Secret:            os.Getenv("APCA_API_SECRET_KEY"),
		Symbols:           tickers,
		RawMessageChannel: rawMessageChannel,
		ErrorsChannel:     errorsChannel,
		Options:           tradeRipperOptions,
	})

	go func() {
		gappersTicker := time.NewTicker(1 * time.Second)

		for range gappersTicker.C {
			tradesWriterLock.RLock()
			gappers := analysis.GetGappers(latestPrices)
			tradesWriterLock.RUnlock()
			server.PublishEvent(map[string]any{
				"type":    "gappers",
				"message": "update",
				"data":    gappers,
			})
		}
	}()

	go func() {
		marketStatusTicker := time.NewTicker(1 * time.Second)

		for range marketStatusTicker.C {
			status := market.GetStatus()
			server.PublishEvent(map[string]any{
				"message": "update",
				"type":    "market_status",
				"data":    status,
			})
		}
	}()

	go func() {
		count, countErr := queries.Count(tradeRipperOptions)

		if countErr != nil {
			logrus.Panicf("Error counting trades: %s", countErr)
		}
		countTicker := time.NewTicker(1 * time.Second)

		for range countTicker.C {
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
						"type":    "trade_count",
						"message": "count update",
						"data": map[string]any{
							"n": count,
						},
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
				for _, trade := range tradeMessages {
					latestPrices[trade["S"].(string)] = trade["p"].(float64)
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
