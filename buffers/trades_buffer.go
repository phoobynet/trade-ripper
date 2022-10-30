package buffers

import (
	"context"
	"fmt"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/queries"
	"github.com/phoobynet/trade-ripper/server"
	"github.com/phoobynet/trade-ripper/utils"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

const BatchSize = 2_000

type TradesBuffer struct {
	sender      *qdb.LineSender
	totalTrades int64
	bufferCount int64
	buffer      []alpaca.TradeRow
	ctx         context.Context
	mu          sync.Mutex
	options     configuration.Options
}

func NewQuestBuffer(options configuration.Options) *TradesBuffer {
	count, countErr := queries.Count(options)

	if countErr != nil {
		panic(countErr)
	}

	questDBAddress := fmt.Sprintf("%s:%d", options.DBHost, options.DBInfluxPort)
	logrus.Infof("Connecting to %s", questDBAddress)

	sender, err := qdb.NewLineSender(context.TODO(), qdb.WithAddress(questDBAddress))

	if err != nil {
		logrus.Fatal("Error creating QuestDB line sender: ", err)
	}

	logrus.Infof("Attempting to connect to %s...CONNECTED", questDBAddress)

	return &TradesBuffer{
		sender:      sender,
		ctx:         context.Background(),
		options:     options,
		totalTrades: count,
	}
}

func (q *TradesBuffer) Start() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		q.flush()
		server.Broadcast(
			server.TradeCountMessage{
				Message: server.Message{
					Type: "tradeCount",
				},
				Count: q.totalTrades,
			})
	}
}

func (q *TradesBuffer) Add(trade alpaca.TradeRow) {
	if strings.HasSuffix(trade.Symbol, "TEST.A") {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.bufferCount += 1
	q.totalTrades += q.bufferCount

	q.buffer = append(q.buffer, trade)
}

func (q *TradesBuffer) flush() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.bufferCount == 0 {
		return
	}

	tradeBatches := utils.Chunk(q.buffer, BatchSize)

	for _, tradeBatch := range tradeBatches {
		for _, trade := range tradeBatch {
			insertErr := q.sender.Table(q.options.Class).Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).At(q.ctx, trade.Timestamp)

			if insertErr != nil {
				logrus.Error("failed to send trade to quest: ", insertErr)
			}
		}
	}

	q.bufferCount = 0
	q.buffer = make([]alpaca.TradeRow, 0)

	err := q.sender.Flush(q.ctx)

	if err != nil {
		logrus.Errorf("error inserting docs: %s", err)
	}
}
