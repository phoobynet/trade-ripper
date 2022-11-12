package alpaca

import (
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model
	ID                           string  `json:"id" gorm:"primaryKey"`
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

func (Asset) TableName() string {
	return "assets"
}

func InitAssets() error {
	db := localdb.Get()
	autoMigrateErr := db.AutoMigrate(&Asset{})

	if autoMigrateErr != nil {
		return autoMigrateErr
	}

	if count, countErr := getAssetCount(); countErr == nil {
		if count == 0 {
			assets, getDataErr := GetData[[]Asset]("/assets", map[string]string{
				"status": "active",
			})

			if getDataErr != nil {
				return getDataErr
			}

			chunks := lo.Chunk(assets, 20)

			for _, chunk := range chunks {
				createErr := db.Model(&Asset{}).Create(&chunk).Error

				if createErr != nil {
					return createErr
				}
			}

		}
	} else {
		return countErr
	}
	return nil
}

func getAssetCount() (int64, error) {
	db := localdb.Get()

	var count int64

	countErr := db.Model(&Asset{}).Count(&count).Error

	if countErr != nil {
		return 0, countErr
	}
	return count, nil
}

func GetAsset(ticker string) (*Asset, error) {
	db := localdb.Get()

	var asset Asset

	findErr := db.Model(&Asset{}).Where("symbol = ?", ticker).First(&asset).Error

	if findErr != nil {
		return nil, findErr
	}

	return &asset, nil
}

func GetAssets(tickers []string) (map[string]Asset, error) {
	db := localdb.Get()

	var assets []Asset

	findErr := db.Model(&Asset{}).Where("symbol IN ?", tickers).Find(&assets).Error

	if findErr != nil {
		return nil, findErr
	}

	assetMap := make(map[string]Asset)

	for _, asset := range assets {
		assetMap[asset.Symbol] = asset
	}

	return assetMap, nil
}
