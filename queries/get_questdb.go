package queries

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func GetQuestDB() (*gorm.DB, error) {
	if questDB == nil {
		return nil, fmt.Errorf("questDB is nil.  You need to call initQuestDB() first")
	}

	return questDB, nil
}
