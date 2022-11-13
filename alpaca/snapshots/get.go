package snapshots

import (
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"strings"
)

func Get(tickers []string) map[string]Snapshot {
	chunks := lo.Chunk[string](tickers, 500)

	snapshots := make(map[string]Snapshot, 0)

	for _, chunk := range chunks {
		symbols := strings.Join(chunk, ",")
		logrus.Infof("Fetching snapshots chunk for %s", symbols)
		response, err := alpaca.GetMarketData[map[string]Snapshot]("/stocks/snapshots", map[string]string{
			"symbols": symbols,
		})

		if err != nil {
			panic(err)
		}

		for ticker, snapshot := range response {
			snapshots[ticker] = snapshot
		}
	}

	return snapshots
}
