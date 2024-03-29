package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/phoobynet/trade-ripper/tradesdb"
	"github.com/sirupsen/logrus"
	"time"
)

// AdaptRawMessageToTrades - converts the raw message from the websocket to a slice Trade maps type objects
func AdaptRawMessageToTrades(rawMessageData []byte) ([]tradesdb.Trade, error) {
	var inputMessages []map[string]any
	var result []tradesdb.Trade

	err := json.Unmarshal(rawMessageData, &inputMessages)

	if err != nil {
		return result, fmt.Errorf("failed to unmarshal raw message data: %w", err)
	}

	for _, message := range inputMessages {
		if t, exists := message["T"]; exists {
			if t == "t" {
				timestampRaw := message["t"].(string)

				timestamp, timestampErr := time.Parse(time.RFC3339Nano, timestampRaw)

				if timestampErr != nil {
					logrus.Errorf("failed to parse timestamp %v", timestampRaw)
					continue
				}

				message["t"] = timestamp.UnixNano()
				result = append(result, message)
			} else if t == "error" {
				logrus.Errorf("alpaca error %v=>%v", message["code"], message["msg"])
			} else if t == "success" {
				logrus.Infof("alpaca success: %v", message["msg"])
			} else if t == "subscription" {
				logrus.Info("alpaca subscription message")
			}
		}
	}

	return result, nil
}
