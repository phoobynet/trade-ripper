package alpaca

import (
	"encoding/json"
	"fmt"
	"github.com/phoobynet/trade-ripper/buffers"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func ConvertToTradeRows(rawMessageData []byte) ([]buffers.TradeRow, error) {
	var inputMessages []map[string]any

	err := json.Unmarshal(rawMessageData, &inputMessages)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw message data: %w", err)
	}

	var tradeRows []buffers.TradeRow

	var tradeRow buffers.TradeRow

	for _, message := range inputMessages {
		if t, exists := message["T"]; exists {
			if t == "t" {
				symbol := message["S"].(string)

				if strings.HasSuffix(symbol, "TEST.A") {
					continue
				}

				timestampRaw := message["t"].(string)

				timestamp, timestampErr := time.Parse(time.RFC3339Nano, timestampRaw)

				if timestampErr != nil {
					logrus.Errorf("failed to parse timestamp %v", timestampRaw)
					continue
				}

				tradeRow = buffers.TradeRow{
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
