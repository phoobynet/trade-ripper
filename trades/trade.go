package trades

type Trade struct {
	Timestamp int64
	Size      float64
	Price     float64
	Symbol    string
	Tks       string
	Base      string
	Quote     string
}
