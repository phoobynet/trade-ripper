package postgres

import (
	"fmt"
	"github.com/phoobynet/trade-ripper/configuration"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(options configuration.Options) {
	dsn := fmt.Sprintf("host=%s user=admin password=quest dbname=qdb port=%d sslmode=disable", options.DBHost, options.DBPostgresPort)

	questDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logrus.Fatal(err)
	}

	db = questDB
}
