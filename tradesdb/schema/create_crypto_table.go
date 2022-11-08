package schema

import (
	_ "embed"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
)

//go:embed create_crypto_table.sql
var createCryptoTableSQL string

func CreateCryptoTable() error {
	db, err := postgres.Get()

	if err != nil {
		return err
	}

	return db.Exec(createCryptoTableSQL).Error
}
