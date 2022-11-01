package buffers

import (
	"context"
	"fmt"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/queries"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const BatchSize = 5_000

type TradesBuffer struct {
	sender     *qdb.LineSender
	bufferLock sync.Mutex
	buffer     [][]byte
	ctx        context.Context
	options    configuration.Options
	tradeCount int64
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
		q.flush()
		tradeCountLog := logrus.WithFields(logrus.Fields{
			"n": q.tradeCount,
		})

		tradeCountLog.Info("count")
	}
}

func (q *TradesBuffer) Add(rawMessage []byte) {
	q.bufferLock.Lock()
	defer q.bufferLock.Unlock()
	q.buffer = append(q.buffer, rawMessage)
}

func (q *TradesBuffer) flush() {
	q.bufferLock.Lock()
	defer q.bufferLock.Unlock()
	var insertErr error
	var tradeRows []alpaca.TradeRow

	for _, rawMessage := range q.buffer {
		rows, conversionErr := alpaca.ConvertToTradeRows(rawMessage)
		if conversionErr != nil {
			logrus.Error(conversionErr)
			continue
		}

		tradeRows = append(tradeRows, rows...)
	}

	tradeBatches := lo.Chunk(tradeRows, BatchSize)

	for _, tradeBatch := range tradeBatches {
		for _, trade := range tradeBatch {
			if trade.Tks == "" {
				insertErr = q.sender.Table(q.options.Class).Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).At(q.ctx, trade.Timestamp)
			} else {
				insertErr = q.sender.Table(q.options.Class).Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).StringColumn("tks", trade.Tks).StringColumn("b", trade.Base).StringColumn("q", trade.Quote).At(q.ctx, trade.Timestamp)
			}

			if insertErr != nil {
				logrus.Error("failed to send trade to quest: ", insertErr)
			}
		}
		err := q.sender.Flush(q.ctx)

		if err != nil {
			logrus.Errorf("error inserting docs: %s", err)
		}
	}

	q.buffer = make([][]byte, 0)
	q.tradeCount += int64(len(tradeRows))
}
