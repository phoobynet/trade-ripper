package localdb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
)

func Get() *gorm.DB {
	if db == nil {
		var dbPath string
		wd, _ := os.Getwd()

		dbPath = filepath.Join(wd, "trade-ripper.db")

		d, dErr := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if dErr != nil {
			panic(dErr)
		}

		db = d
	}

	return db
}
