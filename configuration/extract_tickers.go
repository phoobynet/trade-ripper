package configuration

import (
	"bufio"
	"fmt"
	"github.com/phoobynet/trade-ripper/scrapers"
	"github.com/phoobynet/trade-ripper/utils"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func (o *Options) ExtractTickers() []string {
	tickers := make([]string, 0)

	indexConstituents := make([]scrapers.IndexConstituent, 0)

	if o.Indexes != "" {
		for _, index := range strings.Split(o.Indexes, ",") {
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

	return tickers
}
