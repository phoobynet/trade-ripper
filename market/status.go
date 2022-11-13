package market

import "github.com/phoobynet/trade-ripper/alpaca/calendars"

type Status struct {
	LocalTime  string             `json:"localTime"`
	MarketTime string             `json:"marketTime"`
	Status     string             `json:"status"`
	Current    calendars.Calendar `json:"current"`
	Previous   calendars.Calendar `json:"previous"`
	Next       calendars.Calendar `json:"next"`
}
