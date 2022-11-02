package trades

import (
	_ "embed"
	"github.com/phoobynet/trade-ripper/queries"
)

//go:embed create_us_equity_table.sql
var createUSEquityTableSQL string

func createUSEquityTable() error {
	db, err := queries.GetQuestDB()

	if err != nil {
		return err
	}

	return db.Exec(createUSEquityTableSQL).Error
}
