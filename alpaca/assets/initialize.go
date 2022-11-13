package assets

import (
	"github.com/phoobynet/trade-ripper/alpaca"
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func Initialize() {
	logrus.Infoln("Initializing assets...")
	db := localdb.Get()
	autoMigrateErr := db.AutoMigrate(&Asset{})

	if autoMigrateErr != nil {
		logrus.Fatal(autoMigrateErr)
	}

	if count() == 0 {
		logrus.Infoln("Fetching assets...")
		rawAssets, getDataErr := alpaca.GetTradeData[[]RawAsset]("/assets", map[string]string{
			"status":      "active",
			"asset_class": "us_equity",
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
			createErr := db.Model(&Asset{}).Create(&chunk).Error

			if createErr != nil {
				logrus.Fatal(createErr)
			}
		}
		logrus.Infoln("Initializing assets...DONE")
	} else {
		logrus.Infoln("Assets already initialized")
	}
}
