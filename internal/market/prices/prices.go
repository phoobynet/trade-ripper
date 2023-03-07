package prices

import (
	calendars2 "github.com/phoobynet/trade-ripper/internal/market/calendars"
	"github.com/phoobynet/trade-ripper/internal/market/snapshots"
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
	"time"
)

func CachePreviousClose(tickers []string) {
	db := localdb.Get()
	autoMigrateErr := db.AutoMigrate(&PreviousClose{})

	if autoMigrateErr != nil {
		logrus.Fatal(autoMigrateErr)
	}

	previousCalendar := calendars2.Previous()

	var count int64

	countErr := db.Raw("select count(*) from previous_close where date = ?", previousCalendar.Date).Scan(&count).Error

	if countErr != nil {
		logrus.Fatal(countErr)
	}

	if count == 0 {
		logrus.Info("Caching previous close...")
		tickersSnapshots := snapshots.GetSnapshots(tickers)

		previousCloses := make([]PreviousClose, 0)

		today := time.Now().Format("2006-01-02")

		var price float64
		for ticker, s := range tickersSnapshots {
			if reflect.ValueOf(s).IsZero() || reflect.ValueOf(s.DailyBar).IsZero() || reflect.ValueOf(s.PrevDailyBar).IsZero() {
				continue
			}

			if s.DailyBar.Timestamp == "" {
				continue
			}

			if today == s.DailyBar.Timestamp[:10] {
				price = s.PrevDailyBar.Close
			} else {
				price = s.DailyBar.Close
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

func GetPreviousClosingPrices() (map[string]float64, calendars2.Calendar) {
	db := localdb.Get()

	previousCalendar := calendars2.Previous()

	var previousCloses []PreviousClose

	err := db.Raw("select ticker, price, date from previous_close where date = ?", previousCalendar.Date).Scan(&previousCloses).Error

	if err != nil {
		logrus.Fatal(err)
	}

	previousClosingPrices := make(map[string]float64, 0)

	for _, previousClose := range previousCloses {
		previousClosingPrices[previousClose.Ticker] = previousClose.Price
	}

	return previousClosingPrices, *previousCalendar
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
