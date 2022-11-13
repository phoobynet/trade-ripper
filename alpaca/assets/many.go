package assets

import (
	"github.com/phoobynet/trade-ripper/localdb"
)

func Many(tickers []string) (map[string]Asset, error) {
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
