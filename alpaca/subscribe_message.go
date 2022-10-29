package alpaca

type subscribeMessage struct {
	Action string   `json:"action"`
	Trades []string `json:"trades,omitempty"`
}
