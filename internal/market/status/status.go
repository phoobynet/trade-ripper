package status

import (
	calendars2 "github.com/phoobynet/trade-ripper/internal/market/calendars"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	nextCalendar     *calendars2.Calendar
	previousCalendar *calendars2.Calendar
	currentCalendar  *calendars2.Calendar
	date             time.Time
	location         *time.Location
	marketTimezone   = "America/New_York"
)

func init() {
	location, _ = time.LoadLocation(marketTimezone)
}

func GetMarketStatus() Status {
	localTime := time.Now()
	marketTime := localTime.In(location)

	if localTime.Day() != date.Day() {
		logrus.Debug("New day, refreshing calendars")
		date = localTime
		nextCalendar = calendars2.Next()
		previousCalendar = calendars2.Previous()
		currentCalendar = calendars2.Current()
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

type Status struct {
	LocalTime  time.Time            `json:"localTime"`
	MarketTime time.Time            `json:"marketTime"`
	Status     string               `json:"status"`
	Current    *calendars2.Calendar `json:"current"`
	Previous   *calendars2.Calendar `json:"previous"`
	Next       *calendars2.Calendar `json:"next"`
}
