package market

import (
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var nextCalendar *calendars.Calendar
var previousCalendar *calendars.Calendar
var currentCalendar *calendars.Calendar

var date time.Time

var location *time.Location

func init() {
	location, _ = time.LoadLocation("America/New_York")
}

func GetStatus() Status {
	localTime := time.Now()
	marketTime := localTime.In(location)

	if localTime.Day() != date.Day() {
		logrus.Debug("New day, refreshing calendars")
		date = localTime
		nextCalendar = calendars.Next()
		previousCalendar = calendars.Previous()
		currentCalendar = calendars.Current()
	}

	status := "closed_today"

	if currentCalendar != nil {
		if marketTime.Before(currentCalendar.SessionOpen) {
			status = "opening_later"
		} else if marketTime.Equal(currentCalendar.SessionOpen) || marketTime.Before(currentCalendar.Open) {
			status = "pre_market"
		} else if marketTime.Equal(currentCalendar.Open) || marketTime.Before(currentCalendar.Close) {
			status = "open"
		} else if marketTime.Equal(currentCalendar.Close) || marketTime.Before(currentCalendar.SessionClose) {
			status = "post_market"
		} else if marketTime.Equal(currentCalendar.SessionClose) || marketTime.After(currentCalendar.SessionClose) {
			status = "closed"
		}
	}

	if strings.HasPrefix(status, "closed") {
		currentCalendar = nil
	}

	return Status{
		LocalTime:  localTime,
		MarketTime: marketTime,
		Status:     status,
		Current:    currentCalendar,
		Previous:   previousCalendar,
		Next:       nextCalendar,
	}
}
