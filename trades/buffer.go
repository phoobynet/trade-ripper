package trades

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Buffer struct {
	sender                  *qdb.LineSender
	ctx                     context.Context
	options                 configuration.Options
	write                   func(Trade)
	tradeCount              int64
	tradeBufferPendingCount int64
}

func NewBuffer(options configuration.Options, write func(Trade)) *Buffer {
	return &Buffer{
		ctx:     context.Background(),
		options: options,
		write:   write,
	}
}

func (b *Buffer) Add(rawMessage []byte) {
	trades, convertToTradeErr := convertToTrades(rawMessage)

	if convertToTradeErr != nil {
		logrus.Error(convertToTradeErr)
		return
	}

	for _, trade := range trades {
		b.write(trade)
	}
}

func convertToTrades(rawMessageData []byte) ([]Trade, error) {
	var inputMessages []map[string]any
	var tradeRows []Trade

	err := json.Unmarshal(rawMessageData, &inputMessages)

	if err != nil {
		return tradeRows, fmt.Errorf("failed to unmarshal raw message data: %w", err)
	}

	var tradeRow Trade
	var symbol string

	for _, message := range inputMessages {
		if t, exists := message["T"]; exists {
			if t == "t" {
				symbol = message["S"].(string)

				if strings.HasSuffix(symbol, "TEST.A") {
					continue
				}

				timestampRaw := message["t"].(string)

				timestamp, timestampErr := time.Parse(time.RFC3339Nano, timestampRaw)

				if timestampErr != nil {
					logrus.Errorf("failed to parse timestamp %v", timestampRaw)
					continue
				}

				tradeRow = Trade{
					Symbol:    symbol,
					Size:      message["s"].(float64),
					Price:     message["p"].(float64),
					Timestamp: timestamp.UnixNano(),
				}

				tks, tksExists := message["tks"]

				if tksExists {
					baseQuote := strings.Split(symbol, "/")
					tradeRow.Tks = tks.(string)
					tradeRow.Base = baseQuote[0]
					tradeRow.Quote = baseQuote[1]
				} else {
					tradeRow.Tks = ""
					tks = ""
				}

				tradeRows = append(tradeRows, tradeRow)
			} else if t == "error" {
				logrus.Errorf("alpaca error %v=>%v", message["code"], message["msg"])
			} else if t == "success" {
				logrus.Infof("alpaca success: %v", message["msg"])
			} else if t == "subscription" {
				logrus.Info("alpaca subscription message")
			}
		}
	}

	return tradeRows, nil
}
