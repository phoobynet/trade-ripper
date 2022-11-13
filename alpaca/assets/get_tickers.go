package assets

import "github.com/phoobynet/trade-ripper/localdb"

func GetTickers() []string {
	db := localdb.Get()

	rows, tickersErr := db.Model(&Asset{}).Select("symbol").Distinct().Rows()

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
