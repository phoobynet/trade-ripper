package snapshots

import (
	"github.com/phoobynet/trade-ripper/internal/market"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"strings"
)

func GetSnapshots(tickers []string) map[string]Snapshot {
	chunks := lo.Chunk[string](tickers, 500)

	snapshots := make(map[string]Snapshot, 0)

	for _, chunk := range chunks {
		symbols := strings.Join(chunk, ",")
		logrus.Infof("Fetching snapshots chunk for %s", symbols)
		response, err := market.GetMarketData[map[string]Snapshot]("/stocks/snapshots", map[string]string{
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

type SnapshotTrade struct {
	Timestamp  string   `json:"t"`
	Exchange   string   `json:"x"`
	Price      float64  `json:"p"`
	Size       float64  `json:"s"`
	Conditions []string `json:"c"`
	Index      float64  `json:"i"`
	Tape       string   `json:"z"`
}

type SnapshotBar struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	Timestamp string  `json:"t"`
}

type SnapshotQuote struct {
	AskExchange string   `json:"ax"`
	AskPrice    float64  `json:"ap"`
	AskSize     float64  `json:"as"`
	BidExchange string   `json:"bx"`
	BidPrice    float64  `json:"bp"`
	BidSize     float64  `json:"bs"`
	Timestamp   string   `json:"t"`
	Conditions  []string `json:"c"`
}

type Snapshot struct {
	LatestTrade  SnapshotTrade `json:"latestTrade"`
	LatestQuote  SnapshotQuote `json:"latestQuote"`
	MinuteBar    SnapshotBar   `json:"minuteBar"`
	DailyBar     SnapshotBar   `json:"dailyBar"`
	PrevDailyBar SnapshotBar   `json:"prevDailyBar"`
}
