package queries

import (
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitQuestDB(options configuration.Options) error {
	dsn := fmt.Sprintf("host=%s user=admin password=quest dbname=qdb port=%d sslmode=disable", options.DBHost, options.DBPostgresPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	questDB = db

	return nil
}
