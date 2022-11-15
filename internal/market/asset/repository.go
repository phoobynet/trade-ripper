package asset

import (
	"github.com/phoobynet/trade-ripper/internal/market"
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once

var repositoryInstance *Repository

type Repository struct {
	db *gorm.DB
}

func GetRepositoryInstance() *Repository {
	if repositoryInstance == nil {
		once.Do(func() {
			repositoryInstance = &Repository{
				db: localdb.Get(),
			}
			repositoryInstance.initialize()
		})
	}

	return repositoryInstance
}

func (r *Repository) Tickers() []string {
	rows, tickersErr := r.db.Model(&Asset{}).Select("symbol").Distinct().Rows()

	if tickersErr != nil {
		panic(tickersErr)
	}

	var tickers []string

	for rows.Next() {
		var ticker string
		_ = rows.Scan(&ticker)
		tickers = append(tickers, ticker)
	}

	return tickers
}

func (r *Repository) Many(tickers []string) (map[string]Asset, error) {
	var assets []Asset

	findErr := r.db.Model(&Asset{}).Where("symbol IN ?", tickers).Find(&assets).Error

	if findErr != nil {
		return nil, findErr
	}

	assetMap := make(map[string]Asset)

	for _, asset := range assets {
		assetMap[asset.Symbol] = asset
	}

	return assetMap, nil
}

func (r *Repository) ManySimplified(tickers []string) map[string]SimplifiedAsset {
	assets, err := r.Many(tickers)

	if err != nil {
		panic(err)
	}

	assetMap := make(map[string]SimplifiedAsset)

	for _, asset := range assets {
		assetMap[asset.Symbol] = SimplifiedAsset{
			Symbol:   asset.Symbol,
			Name:     asset.Name,
			Exchange: asset.Exchange,
		}
	}

	return assetMap
}

func (r *Repository) initialize() {
	logrus.Infoln("Initializing assets...")
	autoMigrateErr := r.db.AutoMigrate(&Asset{})

	if autoMigrateErr != nil {
		logrus.Fatal(autoMigrateErr)
	}

	if r.count() == 0 {
		logrus.Infoln("Fetching assets...")
		rawAssets, getDataErr := market.GetTradeData[[]RawAsset]("/assets", map[string]string{
			"status": "active",
		})

		if getDataErr != nil {
			logrus.Fatal(getDataErr)
		}

		assets := lo.Map[RawAsset, Asset](rawAssets, func(rawAsset RawAsset, _ int) Asset {
			return Asset{
				ID:                           rawAsset.ID,
				Symbol:                       rawAsset.Symbol,
				Name:                         rawAsset.Name,
				Exchange:                     rawAsset.Exchange,
				Class:                        rawAsset.Class,
				Status:                       rawAsset.Status,
				Tradable:                     rawAsset.Tradable,
				Fractionable:                 rawAsset.Fractionable,
				Shortable:                    rawAsset.Shortable,
				EasyToBorrow:                 rawAsset.EasyToBorrow,
				Marginable:                   rawAsset.Marginable,
				MinOrderSize:                 rawAsset.MinOrderSize,
				MinTradeIncrement:            rawAsset.MinTradeIncrement,
				PriceIncrement:               rawAsset.PriceIncrement,
				MaintenanceMarginRequirement: rawAsset.MaintenanceMarginRequirement,
			}
		})

		chunks := lo.Chunk(assets, 20)

		logrus.Infoln("Caching assets...")
		for _, chunk := range chunks {
			createErr := r.db.Model(&Asset{}).Create(&chunk).Error

			if createErr != nil {
				logrus.Fatal(createErr)
			}
		}
		logrus.Infoln("Initializing assets...DONE")
	} else {
		logrus.Infoln("Assets already initialized")
	}
}

func (r *Repository) GetOne(ticker string) (*Asset, error) {
	return nil, nil
}

func (r *Repository) count() int64 {
	var c int64

	countErr := r.db.Model(&Asset{}).Count(&c).Error

	if countErr != nil {
		logrus.Fatal(countErr)
	}

	return c
}

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

type SimplifiedAsset struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
}
