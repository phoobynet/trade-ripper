package alpaca

import (
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/server"
)

type TradeReaderConfig struct {
	Key               string
	Secret            string
	Symbols           []string
	ErrorsChannel     chan error
	RawMessageChannel chan []byte
	Options           configuration.Options
	webServer         *server.WebServer
}
