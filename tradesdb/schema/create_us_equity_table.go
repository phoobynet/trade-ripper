package schema

import (
	_ "embed"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
)

//go:embed create_us_equity_table.sql
var createUSEquityTableSQL string

func CreateUSEquityTable() error {
	db, err := postgres.Get()

	if err != nil {
		return err
	}

	return db.Exec(createUSEquityTableSQL).Error
}
