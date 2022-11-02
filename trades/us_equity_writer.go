package trades

import (
	"context"
	"github.com/phoobynet/trade-ripper/configuration"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type USEquityWriter struct {
	sender *qdb.LineSender
	ctx    context.Context
}

func NewUSEquityWriter(options configuration.Options) *USEquityWriter {
	createTableErr := createUSEquityTable()

	if createTableErr != nil {
		logrus.Fatal("Error creating us_equity table: ", createTableErr)
	}

	ctx := context.TODO()
	sender := createSender(ctx, options)

	return &USEquityWriter{
		sender: sender,
		ctx:    ctx,
	}
}

// Write - writes and flushes the trades to QuestDB - recommended to be called when the trades collection reaches between 10 and 1000 objects
func (w *USEquityWriter) Write(trades []Trade) {
	logrus.Infof("Writing %d trades to QuestDB", len(trades))
	chunks := lo.Chunk(trades, 1_000)

	var table *qdb.LineSender

	for _, chunk := range chunks {
		for _, trade := range chunk {
			table = w.sender.Table("us_equity")
			table.Symbol("ticker", trade["S"].(string))
			table.Float64Column("size", trade["s"].(float64))
			table.Float64Column("price", trade["p"].(float64))
			table.At(w.ctx, trade["t"].(int64))
		}
		w.sender.Flush(w.ctx)
	}
}
