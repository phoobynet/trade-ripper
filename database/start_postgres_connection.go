package database

import (
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartPostgresConnection(options configuration.Options) error {
	dsn := fmt.Sprintf("host=%s user=admin password=quest dbname=qdb port=%d sslmode=disable", options.DBHost, options.DBPostgresPort)

	questDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	db = questDB
	return nil
}
