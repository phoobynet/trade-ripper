package alpaca

type TradeRow struct {
	Symbol    string
	Size      float64
	Price     float64
	Tks       string
	Base      string
	Quote     string
	Timestamp int64
}
