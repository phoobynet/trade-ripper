package queries

import "time"

type GetOHLCVOptions struct {
	Ticker   string
	Interval string
	Limit    int
	Date     time.Time
}

type OHLCV struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func GetOHLCV(options GetOHLCVOptions) []OHLCV {
	return []OHLCV{}
}
