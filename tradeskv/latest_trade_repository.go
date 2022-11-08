package tradeskv

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/phoobynet/trade-ripper/tradesdb"
	"strconv"
	"strings"
	"time"
)

// LatestTradeRepository is a repository for the latest trade stored in a key-value data store for performance, alleviating the need for reads from QuestDB to avoid contention.
type LatestTradeRepository struct {
	db      *badger.DB
	options configuration.Options
}

func NewLatestRepository(options configuration.Options) *LatestTradeRepository {
	latestTradesDB, badgerErr := badger.Open(badger.DefaultOptions("latest_trades"))

	if badgerErr != nil {
		panic(badgerErr)
	}

	return &LatestTradeRepository{latestTradesDB, options}
}

func (r *LatestTradeRepository) Update(trades []tradesdb.Trade) error {
	return r.db.Update(func(txn *badger.Txn) error {
		for _, trade := range trades {
			if r.options.Class == "crypto" {
				tks, _ := trade["tks"]
				_ = txn.Set([]byte(trade["S"].(string)), []byte(fmt.Sprintf("%6.4f,%6.4f,%d,%s", trade["s"].(float64), trade["p"].(float64), trade["t"].(int64), tks)))
			} else {
				_ = txn.Set([]byte(trade["S"].(string)), []byte(fmt.Sprintf("%6.4f,%6.4f,%d", trade["s"].(float64), trade["p"].(float64), trade["t"].(int64))))
			}
		}
		return nil
	})
}

func (r *LatestTradeRepository) Get(tickers []string) (map[string]any, error) {
	trades := make(map[string]any)
	err := r.db.View(func(txn *badger.Txn) error {
		for _, ticker := range tickers {
			trade, err := txn.Get([]byte(strings.ToUpper(ticker)))
			if err != nil {
				if err == badger.ErrKeyNotFound {
					continue
				} else {
					return err
				}
			}

			tradeErr := trade.Value(func(val []byte) error {
				tokens := strings.Split(string(val), ",")
				size, sizeErr := strconv.ParseFloat(tokens[0], 64)
				if sizeErr != nil {
					return sizeErr
				}
				price, priceErr := strconv.ParseFloat(tokens[1], 64)

				if priceErr != nil {
					return priceErr
				}

				timestamp, timestampErr := strconv.ParseInt(tokens[2], 10, 64)

				if timestampErr != nil {
					return timestampErr
				}

				if r.options.Class == "crypto" {
					trades[ticker] = map[string]any{
						"size":      size,
						"price":     price,
						"timestamp": time.Unix(0, timestamp),
						"tks":       tokens[3],
					}
				} else {
					trades[ticker] = map[string]any{
						"size":      size,
						"price":     price,
						"timestamp": time.Unix(0, timestamp),
					}
				}

				return nil
			})

			if tradeErr != nil {
				return tradeErr
			}
		}
		return nil
	})

	if err != nil {
		return trades, err
	}

	return nil, err
}

func (r *LatestTradeRepository) GetKeys() ([]string, error) {
	tickers := make([]string, 0)

	err := r.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			tickers = append(tickers, string(k))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tickers, nil
}

func (r *LatestTradeRepository) Close() {
	_ = r.db.Close()
}
