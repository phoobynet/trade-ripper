package assets

func ManySimplified(tickers []string) map[string]SimplifiedAsset {
	assets, err := Many(tickers)

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
