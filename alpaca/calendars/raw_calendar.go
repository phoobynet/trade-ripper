package calendars

type RawCalendar struct {
	Date         string `json:"date"`
	Open         string `json:"open"`
	SessionOpen  string `json:"session_open"`
	Close        string `json:"close"`
	SessionClose string `json:"session_close"`
}
