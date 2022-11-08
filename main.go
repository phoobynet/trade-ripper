package main

import (
	"embed"
	"github.com/alexflint/go-arg"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/loggers"
	"github.com/phoobynet/trade-ripper/scrapers"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/phoobynet/trade-ripper/tradesdb"
	"github.com/phoobynet/trade-ripper/tradesdb/adapters"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"github.com/phoobynet/trade-ripper/tradesdb/queries"
	"github.com/phoobynet/trade-ripper/tradesdb/writers"
	"github.com/phoobynet/trade-ripper/tradeskv"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	//go:embed dist
	dist                  embed.FS
	quitChannel           = make(chan os.Signal, 1)
	options               configuration.Options
	sipReader             *alpaca.TradeReader
	rawMessageChannel     = make(chan []byte, 1_000_000)
	errorsChannel         = make(chan error, 1)
	errorsReceived        = 0
	tradesChannel         = make(chan []tradesdb.Trade, 10_000)
	tradesBuffer          = make([]tradesdb.Trade, 0)
	tradesWriterLock      = sync.Mutex{}
	latestTradeRepository *tradeskv.LatestTradeRepository
	webServer             *server.WebServer
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

	questDBErr := postgres.Start(options)
	if questDBErr != nil {
		panic(questDBErr)
	}

	latestTradeRepository = tradeskv.NewLatestRepository(options)

	defer func(repository *tradeskv.LatestTradeRepository) {
		repository.Close()
	}(latestTradeRepository)

	webServer = server.NewWebServer(options, dist, latestTradeRepository)

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

	symbols := make([]string, 0)

	indexConstituents := make([]scrapers.IndexConstituent, 0)

	if options.Indexes != "" {
		for _, index := range strings.Split(options.Indexes, ",") {
			if index == "sp500" {
				sp500, sp500Err := scrapers.GetSP500()
				if sp500Err != nil {
					logrus.Fatalln(sp500Err)
				}

				indexConstituents = append(indexConstituents, sp500...)
			} else if index == "nasdaq100" {
				nasdaq100, nasdaq100Err := scrapers.GetNASDAQ100()
				if nasdaq100Err != nil {
					logrus.Fatalln(nasdaq100Err)
				}

				indexConstituents = append(indexConstituents, nasdaq100...)
			} else if index == "djia" {
				djia, djiaErr := scrapers.GetDJIA()
				if djiaErr != nil {
					logrus.Fatalln(djiaErr)
				}
				indexConstituents = append(indexConstituents, djia...)
			}
		}

		if len(indexConstituents) > 0 {
			options.Class = "us_equity"
			indexConstituentsSymbols := make([]string, 0)

			for _, ic := range indexConstituents {
				indexConstituentsSymbols = append(indexConstituentsSymbols, ic.Ticker)
			}

			symbols = lo.Uniq[string](indexConstituentsSymbols)
			logrus.Info("Index constituents: ", strings.Join(symbols, ", "))
		} else {
			logrus.Fatalln("No valid market indexes found")
		}
	} else {
		symbols = append(symbols, "*")
	}

	sipReader = alpaca.NewTradeReader(&alpaca.TradeReaderConfig{
		Key:               os.Getenv("APCA_API_KEY_ID"),
		Secret:            os.Getenv("APCA_API_SECRET_KEY"),
		Symbols:           symbols,
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
					tradesBuffer = make([]tradesdb.Trade, 0)
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
}
