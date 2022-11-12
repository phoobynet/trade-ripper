package analysis

import (
	"github.com/phoobynet/trade-ripper/tradesdb/postgres"
	"strconv"
	"strings"
)

type VolumeLeader struct {
	Ticker string `json:"ticker"`
	Volume string `json:"volume"`
	Price  string `json:"price"`
}

const sql = `
		select ticker, sum(size) volume, last (price) price
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

	return volumeLeaders, nil
}
