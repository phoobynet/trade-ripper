package queries

import (
	_ "embed"
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
)

//go:embed count_crypto.sql
var cryptoSQL string

//go:embed count_us_equity.sql
var usEquitySQL string

func Count(options configuration.Options) (int64, error) {
	db, questDBErr := GetQuestDB()

	if questDBErr != nil {
		return 0, questDBErr
	}

	if options.Class == "" {
		return 0, fmt.Errorf("options.Class is empty")
	}

	var count int64

	var sql string
	if options.Class == "crypto" {
		sql = cryptoSQL
	} else {
		sql = usEquitySQL
	}

	queryErr := db.Raw(sql).Scan(&count).Error

	if queryErr != nil {
		return 0, queryErr
	}

	return count, nil
}
