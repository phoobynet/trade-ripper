package calendars

import (
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/sirupsen/logrus"
	"time"
)

const selectPreviousCalendarSQL = `
SELECT * FROM calendars where date < ? ORDER BY date DESC LIMIT 1
`

func Previous() *Calendar {
	date := time.Now().Format("2006-01-02")

	db := localdb.Get()

	var calendar Calendar

	scanErr := db.Raw(selectPreviousCalendarSQL, date).Scan(&calendar).Error

	if scanErr != nil {
		logrus.Fatal(scanErr)
	}

	return &calendar
}
