package trades

type TradeWriter interface {
	Write(trades []Trade)
}
