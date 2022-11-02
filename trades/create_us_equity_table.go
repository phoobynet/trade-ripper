package trades

import (
	_ "embed"
	"github.com/phoobynet/trade-ripper/database"
)

//go:embed create_us_equity_table.sql
var createUSEquityTableSQL string

func createUSEquityTable() error {
	db, err := database.GetPostgresConnection()

	if err != nil {
		return err
	}

	return db.Exec(createUSEquityTableSQL).Error
}
