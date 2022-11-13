package analysis

import (
	"github.com/phoobynet/trade-ripper/alpaca/assets"
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

type VolumeLeader struct {
	Ticker   string  `json:"ticker"`
	Volume   float64 `json:"volume"`
	Price    float64 `json:"price"`
	Name     string  `json:"name"`
	Exchange string  `json:"exchange"`
}

const sql = `
		select ticker, sum(size) volume, round_half_even(last (price), 2) price
    from
        us_equity
    where
        timestamp in ':date'
    group by
        ticker
    order by
        volume desc
        limit :limit
	`

func GetVolumeLeaders(date string, limit int) ([]VolumeLeader, error) {
	db, dbErr := postgres.Get()

	if dbErr != nil {
		return nil, dbErr
	}

	finalSQL := strings.Replace(sql, ":date", date, 1)
	finalSQL = strings.Replace(finalSQL, ":limit", strconv.Itoa(limit), 1)

	var volumeLeaders []VolumeLeader

	scanErr := db.Raw(finalSQL).Scan(&volumeLeaders).Error

	if scanErr != nil {
		return nil, scanErr
	}

	tickers := lo.Map[VolumeLeader, string](volumeLeaders, func(volumeLeader VolumeLeader, _ int) string {
		return volumeLeader.Ticker
	})

	tickerAssets := assets.ManySimplified(tickers)

	volumeLeaders = lo.Map[VolumeLeader, VolumeLeader](volumeLeaders, func(volumeLeader VolumeLeader, _ int) VolumeLeader {
		if asset, ok := tickerAssets[volumeLeader.Ticker]; ok {
			return VolumeLeader{
				Ticker:   volumeLeader.Ticker,
				Volume:   volumeLeader.Volume,
				Price:    volumeLeader.Price,
				Name:     asset.Name,
				Exchange: asset.Exchange,
			}
		} else {
			return volumeLeader
		}
	})

	return volumeLeaders, nil
}
