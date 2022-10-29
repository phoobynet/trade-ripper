package alpaca

import "net/url"

type SIPReaderConfig struct {
	SocketURL         *url.URL // optional - if you don't supply it, the default value will be used
	Key               string
	Secret            string
	Trades            []string
	Quotes            []string
	Bars              []string
	ErrorsChannel     chan error
	RawMessageChannel chan []byte
}
