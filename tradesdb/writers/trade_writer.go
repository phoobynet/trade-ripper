package writers

import "github.com/phoobynet/trade-ripper/tradesdb"

type TradeWriter interface {
	Write(trades []tradesdb.Trade)
}
