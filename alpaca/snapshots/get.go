package snapshots

import (
	"github.com/phoobynet/trade-ripper/alpaca"
	"strings"
)

func Get(tickers []string) map[string]Snapshot {
	response, err := alpaca.GetMarketData[map[string]Snapshot]("snapshots", map[string]string{
		"symbols": strings.Join(tickers, ","),
	})

	if err != nil {
		panic(err)
	}

	return response
}
