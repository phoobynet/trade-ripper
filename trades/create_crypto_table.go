package trades

import (
	_ "embed"
	"github.com/phoobynet/trade-ripper/database"
)

//go:embed create_crypto_table.sql
var createCryptoTableSQL string

func createCryptoTable() error {
	db, err := database.GetPostgresConnection()

	if err != nil {
		return err
	}

	return db.Exec(createCryptoTableSQL).Error
}
