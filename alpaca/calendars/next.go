package calendars

import (
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/sirupsen/logrus"
	"time"
)

const selectNextCalendarSQL = `
SELECT * FROM calendars where date > ? ORDER BY date ASC LIMIT 1
`

func Next() Calendar {
	date := time.Now().Format("2006-01-02")

	db := localdb.Get()

	var calendar Calendar

	scanErr := db.Raw(selectNextCalendarSQL, date).Scan(&calendar).Error

	if scanErr != nil {
		logrus.Fatal(scanErr)
	}

	return calendar
}
