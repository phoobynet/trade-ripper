package queries

import (
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
)

func Count(options configuration.Options) (int64, error) {
	db, questDBErr := GetQuestDB()

	if questDBErr != nil {
		return 0, questDBErr
	}

	if options.Class == "" {
		return 0, fmt.Errorf("options.Class is empty")
	}

	var count int64

	queryErr := db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE t >= timestamp_floor('d',  now())", options.Class)).Scan(&count).Error

	if queryErr != nil {
		return 0, queryErr
	}

	return count, nil
}
