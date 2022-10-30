package queries

import _ "embed"

//go:embed create_crypto_table.sql
var createCryptoTableSQL string

func CreateCryptoTable() error {
	db, err := GetQuestDB()

	if err != nil {
		return err
	}

	return db.Exec(createCryptoTableSQL).Error
}
