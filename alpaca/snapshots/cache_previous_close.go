package snapshots

import (
	"github.com/phoobynet/trade-ripper/alpaca/assets"
	"github.com/phoobynet/trade-ripper/alpaca/calendars"
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func CachePreviousClose() {
	db := localdb.Get()

	previousCalendar := calendars.Previous()

	var count int64

	countErr := db.Raw("select count(*) from previous_close where date = ?", previousCalendar.Date).Scan(&count).Error

	if countErr != nil {
		logrus.Fatal(countErr)
	}

	if count == 0 {
		logrus.Info("Caching previous close...")
		tickers := assets.GetTickers()
		snapshots := Get(tickers)

		previousCloses := make([]PreviousClose, 0)

		today := time.Now().Format("2006-01-02")

		var price float64
		for ticker, s := range snapshots {
			if today == s.DailyBar.Timestamp[:10] {
				price = s.PrevDailyBar.Close
			} else {
				price = s.DailyBar.Open
			}

			previousCloses = append(previousCloses, PreviousClose{
				Date:   previousCalendar.Date,
				Ticker: ticker,
				Price:  price,
			})
		}

		chunks := lo.Chunk(previousCloses, 50)

		for _, chunk := range chunks {
			chunkErr := db.Model(&PreviousClose{}).Create(chunk).Error

			if chunkErr != nil {
				logrus.Fatal(chunkErr)
			}
		}

		logrus.Info("Cached previous close...DONE")
	}
}

type PreviousClose struct {
	gorm.Model
	Ticker string  `gorm:"index"`
	Date   string  `gorm:"index"`
	Price  float64 `json:"price"`
}

func (PreviousClose) TableName() string {
	return "previous_close"
}
