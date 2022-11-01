package trades

type TradeWriter interface {
	Writer(trade Trade)
}
