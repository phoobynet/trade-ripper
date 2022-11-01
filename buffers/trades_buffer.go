package buffers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/queries"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

const ()

type TradesBuffer struct {
	sender                  *qdb.LineSender
	tradeBufferLock         sync.Mutex
	ctx                     context.Context
	options                 configuration.Options
	tradeCount              int64
	tradeBufferPendingCount int64
}

func NewQuestBuffer(options configuration.Options) *TradesBuffer {
	questDBAddress := fmt.Sprintf("%s:%d", options.DBHost, options.DBInfluxPort)
	logrus.Infof("Connecting to %s", questDBAddress)

	sender, err := qdb.NewLineSender(context.TODO(), qdb.WithAddress(questDBAddress))

	if err != nil {
		logrus.Fatal("Error creating QuestDB line sender: ", err)
	}

	logrus.Infof("Attempting to connect to %s...CONNECTED", questDBAddress)

	return &TradesBuffer{
		sender:  sender,
		ctx:     context.Background(),
		options: options,
	}
}

func (q *TradesBuffer) Start() {
	tradeCount, tradeCountErr := queries.Count(q.options)

	if tradeCountErr != nil {
		panic(tradeCountErr)
	}
	q.tradeCount = tradeCount

	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		q.flush(true)
		tradeCountLog := logrus.WithFields(logrus.Fields{
			"n": q.tradeCount,
		})

		tradeCountLog.Info("count")
	}
}

func (q *TradesBuffer) Add(rawMessage []byte) {
	q.tradeBufferLock.Lock()
	defer q.tradeBufferLock.Unlock()
	tradeRows, convertToTradeErr := convertToTrades(rawMessage)

	if convertToTradeErr != nil {
		logrus.Error(convertToTradeErr)
		return
	}

	var insertErr error

	for _, trade := range tradeRows {
		if trade.Tks == "" {
			insertErr = q.sender.Table(q.options.Class).Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).At(q.ctx, trade.Timestamp)
		} else {
			insertErr = q.sender.Table(q.options.Class).Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).StringColumn("tks", trade.Tks).StringColumn("b", trade.Base).StringColumn("q", trade.Quote).At(q.ctx, trade.Timestamp)
		}

		if insertErr != nil {
			logrus.Error("failed to send trade to quest: ", insertErr)
		}

		q.tradeBufferPendingCount++

		if q.tradeBufferPendingCount >= 1000 {
			err := q.sender.Flush(q.ctx)

			if err != nil {
				logrus.Errorf("error inserting docs: %s", err)
			}
			q.tradeCount += q.tradeBufferPendingCount
			q.tradeBufferPendingCount = 0
		}
	}
}

func (q *TradesBuffer) flush(forceFlush bool) {
	if q.tradeBufferPendingCount >= 1_000 || (forceFlush && q.tradeBufferPendingCount > 0) {
		err := q.sender.Flush(q.ctx)

		if err != nil {
			logrus.Errorf("error inserting docs: %s", err)
		}

		q.tradeCount += q.tradeBufferPendingCount
		q.tradeBufferPendingCount = 0
	}
}

func convertToTrades(rawMessageData []byte) ([]TradeRow, error) {
	var inputMessages []map[string]any
	var tradeRows []TradeRow

	err := json.Unmarshal(rawMessageData, &inputMessages)

	if err != nil {
		return tradeRows, fmt.Errorf("failed to unmarshal raw message data: %w", err)
	}

	var tradeRow TradeRow

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

				tradeRow = TradeRow{
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
