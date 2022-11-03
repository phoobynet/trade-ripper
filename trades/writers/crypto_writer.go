package writers

import (
	"context"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/trades"
	"github.com/phoobynet/trade-ripper/trades/schema"
	qdb "github.com/questdb/go-questdb-client"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type CryptoWriter struct {
	sender *qdb.LineSender
	ctx    context.Context
}

func NewCryptoWriter(options configuration.Options) *CryptoWriter {
	createTableErr := schema.CreateCryptoTable()

	if createTableErr != nil {
		logrus.Fatal("Error creating crypto table: ", createTableErr)
	}

	ctx := context.TODO()
	sender := trades.CreateSender(ctx, options)

	return &CryptoWriter{
		sender: sender,
		ctx:    ctx,
	}
}

// Write - writes and flushes the trades to QuestDB - recommended to be called when the trades reaches between 10 and 1000 objects
func (w *CryptoWriter) Write(trades []trades.Trade) {
	chunks := lo.Chunk(trades, 1_000)

	var table *qdb.LineSender
	var insertErr error
	for _, chunk := range chunks {
		for _, trade := range chunk {
			table = w.sender.Table("crypto")
			table.Symbol("pair", trade["S"].(string))
			table.Float64Column("size", trade["s"].(float64))
			table.Float64Column("price", trade["p"].(float64))
			table.StringColumn("tks", trade["tks"].(string))
			insertErr = table.At(w.ctx, trade["t"].(int64))

			if insertErr != nil {
				logrus.Error("Error inserting crypto trade: ", insertErr)
			}
		}
		w.sender.Flush(w.ctx)
	}

	logrus.Infof("Inserted %d trades into 'crypto'", len(trades))
}
