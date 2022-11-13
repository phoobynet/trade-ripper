package snapshots

type SnapshotTrade struct {
	Timestamp  string   `json:"t"`
	Exchange   string   `json:"x"`
	Price      string   `json:"p"`
	Size       string   `json:"s"`
	Conditions []string `json:"c"`
	Index      float64  `json:"i"`
	Tape       string   `json:"z"`
}

type SnapshotBar struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	Timestamp string  `json:"t"`
}

type SnapshotQuote struct {
	AskExchange string   `json:"ax"`
	AskPrice    float64  `json:"ap"`
	AskSize     float64  `json:"as"`
	BidExchange string   `json:"bx"`
	BidPrice    float64  `json:"bp"`
	BidSize     float64  `json:"bs"`
	Timestamp   string   `json:"t"`
	Conditions  []string `json:"c"`
}

type Snapshot struct {
	LatestTrade  SnapshotTrade `json:"latestTrade"`
	LatestQuote  SnapshotQuote `json:"latestQuote"`
	MinuteBar    SnapshotBar   `json:"minuteBar"`
	DailyBar     SnapshotBar   `json:"dailyBar"`
	PrevDailyBar SnapshotBar   `json:"prevDailyBar"`
}
