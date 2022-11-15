package calendars

import (
	"fmt"
	"github.com/phoobynet/trade-ripper/internal/market"
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const timeZone = "America/New_York"

var newYork *time.Location

func init() {
	newYork, _ = time.LoadLocation("America/New_York")
}

func convertToNewYorkTime(date string, t string) time.Time {
	defer func() {
		if r := recover(); r != nil {
			logrus.Fatalf("Failed to parse the date: %s, time: %s...panic", date, t)
			panic("convertToNewYorkTime failed")
		}
	}()
	var hour string
	var minute string
	if strings.Contains(t, ":") {
		tParts := strings.Split(t, ":")

		hour = tParts[0]
		minute = tParts[1]
	} else {
		hour = t[:2]
		minute = t[2:]
	}

	timeInLocation, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s:%s", date, hour, minute), newYork)

	if err != nil {
		panic(err)
	}

	return timeInLocation
}

func convertCalendar(rawCalendar RawCalendar) Calendar {
	return Calendar{
		Date:         rawCalendar.Date,
		SessionOpen:  convertToNewYorkTime(rawCalendar.Date, rawCalendar.SessionOpen),
		Open:         convertToNewYorkTime(rawCalendar.Date, rawCalendar.Open),
		Close:        convertToNewYorkTime(rawCalendar.Date, rawCalendar.Close),
		SessionClose: convertToNewYorkTime(rawCalendar.Date, rawCalendar.SessionClose),
	}
}

func convertCalendars(rawCalendars []RawCalendar) []Calendar {
	return lo.Map[RawCalendar, Calendar](rawCalendars, func(rawCalendar RawCalendar, _ int) Calendar {
		return convertCalendar(rawCalendar)
	})
}

func Initialize() {
	logrus.Infoln("Initializing calendars...")
	db := localdb.Get()
	autoMigrateErr := db.AutoMigrate(&Calendar{})

	if autoMigrateErr != nil {
		logrus.Fatal(autoMigrateErr)
	}

	if count() == 0 {
		logrus.Infoln("Fetching calendars...")
		t := time.Now()
		start := t.AddDate(-1, 0, 0).Format("2006-01-02")
		end := t.AddDate(1, 0, 0).Format("2006-01-02")

		rawCalendars, getDataErr := market.GetTradeData[[]RawCalendar]("/calendar", map[string]string{
			"start": start,
			"end":   end,
		})

		if getDataErr != nil {
			logrus.Fatal(getDataErr)
		}

		calendars := convertCalendars(rawCalendars)

		chunks := lo.Chunk(calendars, 20)

		logrus.Infoln("Caching calendars...")
		for _, chunk := range chunks {
			createErr := db.Model(&Calendar{}).Create(&chunk).Error

			if createErr != nil {
				logrus.Fatal(createErr)
			}
		}

		logrus.Infoln("Initializing calendars...DONE")
	} else {
		logrus.Infoln("Calendars already initialized")
	}
}
