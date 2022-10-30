package queries

import _ "embed"

//go:embed create_us_equity_table.sql
var createUSEquityTableSQL string

func CreateUSEquityTable() error {
	db, err := GetQuestDB()

	if err != nil {
		return err
	}

	return db.Exec(createUSEquityTableSQL).Error
}
