package writers

import "github.com/phoobynet/trade-ripper/trades"

type TradeWriter interface {
	Write(trades []trades.Trade)
}
