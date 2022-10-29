package alpaca

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func ConvertToTradeRows(rawMessageData []byte) ([]TradeRow, error) {
	var inputMessages []map[string]any

	err := json.Unmarshal(rawMessageData, &inputMessages)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw message data: %w", err)
	}

	var tradeRows []TradeRow

	for _, message := range inputMessages {
		if t, exists := message["T"]; exists {
			if t == "t" {
				symbol := message["S"].(string)

				if strings.HasSuffix(symbol, "TEST.A") {
					continue
				}

				timestampRaw := message["T"].(string)

				timestamp, timestampErr := time.Parse(time.RFC3339Nano, timestampRaw)

				if timestampErr != nil {
					logrus.Errorf("failed to parse timestamp %v", timestampRaw)
					continue
				}

				tradeRows = append(tradeRows, TradeRow{
					Symbol:    symbol,
					Size:      message["s"].(float64),
					Price:     message["p"].(float64),
					Timestamp: timestamp.UnixNano(),
				})
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
