package buffers

import (
	"context"
	"fmt"
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/configuration"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

const AutoFlushThreshold = 2_000

type QuestTradesBuffer struct {
	sender      *qdb.LineSender
	totalTrades int
	bufferCount int
	buffer      []alpaca.TradeRow
	ctx         context.Context
	mu          sync.Mutex
}

func NewQuestBuffer(options configuration.Options) *QuestTradesBuffer {
	questDBAddress := fmt.Sprintf("%s", options.QuestDBURI)
	logrus.Infof("Attempting to connect to %s", questDBAddress)

	sender, err := qdb.NewLineSender(context.TODO(), qdb.WithAddress(questDBAddress))

	if err != nil {
		logrus.Fatal("Error creating QuestDB line sender: ", err)
	}

	logrus.Infof("Attempting to connect to %s...CONNECTED", questDBAddress)

	return &QuestTradesBuffer{
		sender: sender,
		ctx:    context.Background(),
	}
}

func (q *QuestTradesBuffer) Start() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		q.flush()
	}
}

func (q *QuestTradesBuffer) Add(trade alpaca.TradeRow) {
	if strings.HasSuffix(trade.Symbol, "TEST.A") {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.bufferCount += 1
	q.totalTrades += q.bufferCount

	insertErr := q.sender.Table("trades").Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).At(q.ctx, trade.Timestamp)

	if insertErr != nil {
		logrus.Error("failed to send trade to quest: ", insertErr)
	}

	if q.bufferCount > AutoFlushThreshold {
		q.flush()
	}
}

func (q *QuestTradesBuffer) flush() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.bufferCount > 0 {
		err := q.sender.Flush(q.ctx)
		if err != nil {
			logrus.Errorf("error inserting docs: %s", err)
		}
		q.bufferCount = 0
	}
}
