package configuration

import (
	"bufio"
	"fmt"
	"github.com/alexflint/go-arg"
	djia "github.com/phoobynet/djia-constituent-scraper"
	scrapers2 "github.com/phoobynet/trade-ripper/internal/indexes"
	"github.com/phoobynet/trade-ripper/utils"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Options struct {
	DBHost         string `arg:"required,-h,--host" help:"The questdb post e.g. my.questdb.db"`
	DBInfluxPort   int    `arg:"-i,--influx" help:"The questdb influx port e.g. 9009 (default)"`
	DBPostgresPort int    `arg:"-p,--postgres" help:"The questdb postgres port e.g. 8812 (default)"`
	Class          string `arg:"required,-c,--class" help:"The class to subscribe to, either crypto or us_equity"`
	WebServerPort  int    `arg:"-w,--webserver" help:"The webserver port e.g. 3000 (default)"`
	Indexes        string `arg:"--indexes" help:"example: sp500,nasdaq100,djia - Currently only sp500, nasdaq100, djia are supported.  Limits the data to the indexes specified"`
	TickersFile    string `arg:"--tickersfile" help:"example: tickers.txt - A file containing a list of tickers to subscribe to, seperated by newlines.  Can be used in conjunction with --indexes"`
	Tickers        string `arg:"--tickers" help:"example: AAPL,MSFT,GOOG - A comma seperated list of tickers to subscribe to.  Can be used in conjunction with --indexes and --tickersfile"`
}

var tradeRipperOptions Options

func init() {
	arg.MustParse(&tradeRipperOptions)

	if tradeRipperOptions.DBInfluxPort == 0 {
		tradeRipperOptions.DBInfluxPort = 9009
	}

	if tradeRipperOptions.DBPostgresPort == 0 {
		tradeRipperOptions.DBPostgresPort = 8812
	}

	if tradeRipperOptions.WebServerPort == 0 {
		tradeRipperOptions.WebServerPort = 3000
	}

	tradeRipperOptions.extractTickers()
}

func GetTradeRipperOptions() Options {
	return tradeRipperOptions
}

func (o *Options) extractTickers() {
	tickers := make([]string, 0)

	indexConstituents := make([]scrapers2.IndexConstituent, 0)

	if o.Indexes != "" {
		for _, index := range strings.Split(o.Indexes, ",") {
			if index == "sp500" {
				sp500, sp500Err := scrapers2.ScrapeSP500()
				if sp500Err != nil {
					logrus.Fatalln(sp500Err)
				}

				indexConstituents = append(indexConstituents, sp500...)
			} else if index == "nasdaq100" {
				nasdaq100, nasdaq100Err := scrapers2.ScrapeNASDAQ100()
				if nasdaq100Err != nil {
					logrus.Fatalln(nasdaq100Err)
				}

				indexConstituents = append(indexConstituents, nasdaq100...)
			} else if index == "djia" {
				djiaConstituents, djiaErr := djia.ScrapeDJIA()
				if djiaErr != nil {
					logrus.Fatalln(djiaErr)
				}
				for _, djiaConstituent := range djiaConstituents {
					indexConstituents = append(indexConstituents, scrapers2.IndexConstituent{
						Ticker:  djiaConstituent.Ticker,
						Company: djiaConstituent.Company,
					})
				}
			}
		}

		if len(indexConstituents) > 0 {
			o.Class = "us_equity"
			indexConstituentsSymbols := make([]string, 0)

			for _, ic := range indexConstituents {
				indexConstituentsSymbols = append(indexConstituentsSymbols, ic.Ticker)
			}

			tickers = lo.Uniq[string](indexConstituentsSymbols)
			logrus.Info("Index constituents: ", strings.Join(tickers, ", "))
		} else {
			logrus.Fatalln("No valid market indexes found")
		}
	}

	if o.TickersFile != "" {
		if _, err := os.Stat(o.TickersFile); os.IsNotExist(err) {
			logrus.Fatalln("Tickers file in --tickers arg not found")
		}

		f, tickersFileErr := os.OpenFile(o.TickersFile, os.O_RDONLY, 0644)

		if tickersFileErr != nil {
			logrus.Fatalln(fmt.Errorf("failed to open tickers file: %w", tickersFileErr))
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		tickersFileScanner := bufio.NewScanner(f)

		for tickersFileScanner.Scan() {
			ticker := utils.CleanTicker(tickersFileScanner.Text())

			if ticker != "" {
				continue
			}

			tickers = append(tickers, ticker)
		}
	}

	if o.Tickers != "" {
		lo.ForEach[string](strings.Split(o.Tickers, ","), func(ticker string, _ int) {
			cleanedTicker := utils.CleanTicker(ticker)
			if cleanedTicker != "" {
				tickers = append(tickers, cleanedTicker)
			}
		})
	}

	if len(tickers) == 0 {
		tickers = append(tickers, "*")
	}

	o.Tickers = strings.Join(tickers, ",")
}
