package trades

import (
	_ "embed"
	"github.com/phoobynet/trade-ripper/queries"
)

//go:embed create_crypto_table.sql
var createCryptoTableSQL string

func createCryptoTable() error {
	db, err := queries.GetQuestDB()

	if err != nil {
		return err
	}

	return db.Exec(createCryptoTableSQL).Error
}
