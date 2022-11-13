package calendars

import (
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/sirupsen/logrus"
	"time"
)

const selectCurrentCalendarSQL = `SELECT * FROM calendars WHERE date = ? limit 1`

func Current() *Calendar {
	date := time.Now().Format("2006-01-02")

	db := localdb.Get()

	var calendar *Calendar

	scanErr := db.Raw(selectCurrentCalendarSQL, date).Scan(calendar).Error

	if scanErr != nil {
		logrus.Fatal(scanErr)
	}

	return calendar
}
