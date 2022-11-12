package alpaca

import (
	"github.com/phoobynet/trade-ripper/configuration"
)

type TradeReaderConfig struct {
	Key               string
	Secret            string
	Symbols           []string
	ErrorsChannel     chan error
	RawMessageChannel chan []byte
	Options           configuration.Options
}
