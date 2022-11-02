package database

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func GetPostgresConnection() (*gorm.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("questDB is nil.  You need to call initQuestDB() first")
	}

	return db, nil
}
