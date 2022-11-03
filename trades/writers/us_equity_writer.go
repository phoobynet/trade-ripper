package writers

import (
	"context"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/trades"
	"github.com/phoobynet/trade-ripper/trades/schema"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"strings"
)

type USEquityWriter struct {
	sender *qdb.LineSender
	ctx    context.Context
}

func NewUSEquityWriter(options configuration.Options) *USEquityWriter {
	createTableErr := schema.CreateUSEquityTable()

	if createTableErr != nil {
		logrus.Fatal("Error creating us_equity table: ", createTableErr)
	}

	ctx := context.TODO()
	sender := trades.CreateSender(ctx, options)

	return &USEquityWriter{
		sender: sender,
		ctx:    ctx,
	}
}

// Write - writes and flushes the trades to QuestDB - recommended to be called when the trades collection reaches between 10 and 1000 objects
func (w *USEquityWriter) Write(trades []trades.Trade) {
	chunks := lo.Chunk(trades, 1_000)

	var table *qdb.LineSender

	var insertErr error
	var ticker string
	negate := 0
	for _, chunk := range chunks {
		for _, trade := range chunk {
			ticker = trade["S"].(string)
			if strings.HasSuffix(ticker, "TEST.A") {
				negate++
				continue
			}
			table = w.sender.Table("us_equity")
			table.Symbol("ticker", ticker)
			table.Float64Column("size", trade["s"].(float64))
			table.Float64Column("price", trade["p"].(float64))
			insertErr = table.At(w.ctx, trade["t"].(int64))

			if insertErr != nil {
				logrus.Error("Error inserting us_equity trade: ", insertErr)
			}
		}
		flushErr := w.sender.Flush(w.ctx)

		if flushErr != nil {
			logrus.Panicf("Error flushing us_equity trades: %s", flushErr)
		}
	}

	logrus.Infof("Inserted %d trades into 'us_equity'", len(trades)-negate)
}
