package assets

type RawAsset struct {
	ID                           string  `json:"id"`
	Symbol                       string  `json:"symbol" gorm:"index:idx_symbol,unique,sort:asc"`
	Name                         string  `json:"name" gorm:"index:idx_name,sort:asc,not null"`
	Exchange                     string  `json:"exchange" gorm:"not null"`
	Class                        string  `json:"class" gorm:"not null"`
	Status                       string  `json:"status" gorm:"not null"`
	Tradable                     bool    `json:"tradable"`
	Fractionable                 bool    `json:"fractionable"`
	Shortable                    bool    `json:"shortable"`
	EasyToBorrow                 bool    `json:"easy_to_borrow"`
	Marginable                   bool    `json:"marginable"`
	MinOrderSize                 string  `json:"min_order_size,omitempty"`
	MinTradeIncrement            string  `json:"min_trade_increment,omitempty"`
	PriceIncrement               string  `json:"price_increment,omitempty"`
	MaintenanceMarginRequirement float64 `json:"maintenance_margin_requirement,omitempty"`
}
