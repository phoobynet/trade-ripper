package trades

import (
	"context"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/sirupsen/logrus"
	"sync"
)

type Writer struct {
	sender       *qdb.LineSender
	flushLock    sync.Mutex
	ctx          context.Context
	pendingCount int
	options      configuration.Options
}

func NewWriter(options configuration.Options) *Writer {
	questDBAddress := fmt.Sprintf("%s:%d", options.DBHost, options.DBInfluxPort)
	logrus.Infof("Connecting to %s", questDBAddress)

	sender, err := qdb.NewLineSender(context.TODO(), qdb.WithAddress(questDBAddress))

	if err != nil {
		logrus.Fatal("Error creating QuestDB line sender: ", err)
	}

	logrus.Infof("Attempting to connect to %s...CONNECTED", questDBAddress)
	return &Writer{
		sender:  sender,
		ctx:     context.TODO(),
		options: options,
	}
}

func (w *Writer) Write(trades []Trade) {
	w.flushLock.Lock()
	defer w.flushLock.Unlock()

	var insertErr error

	if trade.Tks == "" {
		insertErr = w.sender.Table(w.options.Class).Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).At(w.ctx, trade.Timestamp)
	} else {
		insertErr = w.sender.Table(w.options.Class).Symbol("sy", trade.Symbol).Float64Column("s", trade.Size).Float64Column("p", trade.Price).StringColumn("tks", trade.Tks).StringColumn("b", trade.Base).StringColumn("q", trade.Quote).At(w.ctx, trade.Timestamp)
	}

	if insertErr != nil {
		logrus.Error("failed to send trade to quest: ", insertErr)
	}

	w.pendingCount++
	w.flush(false)
}

func (w *Writer) flush(forceFlush bool) {
	w.flushLock.Lock()
	defer w.flushLock.Unlock()

	if w.pendingCount >= 1_000 || (forceFlush && w.pendingCount > 0) {
		err := w.sender.Flush(w.ctx)

		if err != nil {
			panic(fmt.Errorf("error inserting docs: %w", err))
		}

		w.pendingCount = 0
		fmt.Println("flushed")
	}
}
