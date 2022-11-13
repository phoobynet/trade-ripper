package assets

import "gorm.io/gorm"

type Asset struct {
	gorm.Model                   `json:"-"`
	ID                           string  `json:"-" gorm:"primaryKey"`
	Symbol                       string  `json:"symbol" gorm:"index:idx_symbol,unique,sort:asc"`
	Name                         string  `json:"name" gorm:"index:idx_name,sort:asc,not null"`
	Exchange                     string  `json:"exchange" gorm:"not null"`
	Class                        string  `json:"class" gorm:"not null"`
	Status                       string  `json:"status" gorm:"not null"`
	Tradable                     bool    `json:"tradable"`
	Fractionable                 bool    `json:"fractionable"`
	Shortable                    bool    `json:"shortable"`
	EasyToBorrow                 bool    `json:"easyToBorrow"`
	Marginable                   bool    `json:"marginable"`
	MinOrderSize                 string  `json:"minOrderSize,omitempty"`
	MinTradeIncrement            string  `json:"minTradeIncrement,omitempty"`
	PriceIncrement               string  `json:"priceIncrement,omitempty"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement,omitempty"`
}

func (Asset) TableName() string {
	return "assets"
}
