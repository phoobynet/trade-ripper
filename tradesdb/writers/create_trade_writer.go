package writers

import (
	"github.com/phoobynet/trade-ripper/internal/configuration"
)

func CreateTradeWriter(options configuration.Options) TradeWriter {
	var tradeWriter TradeWriter
	if options.Class == "crypto" {
		tradeWriter = NewCryptoWriter(options)
	} else {
		tradeWriter = NewUSEquityWriter(options)
	}

	return tradeWriter
}
