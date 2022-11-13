package market

import (
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"time"
)

type Status struct {
	LocalTime  time.Time           `json:"localTime"`
	MarketTime time.Time           `json:"marketTime"`
	Status     string              `json:"status"`
	Current    *calendars.Calendar `json:"current"`
	Previous   *calendars.Calendar `json:"previous"`
	Next       *calendars.Calendar `json:"next"`
}
