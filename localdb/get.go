package localdb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Get() *gorm.DB {
	if db == nil {
		d, dErr := gorm.Open(sqlite.Open("local.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if dErr != nil {
			panic(dErr)
		}

		db = d
	}

	return db
}
