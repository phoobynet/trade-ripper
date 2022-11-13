package assets

import (
	"github.com/phoobynet/trade-ripper/localdb"
	"github.com/sirupsen/logrus"
)

func count() int64 {
	db := localdb.Get()

	var c int64

	countErr := db.Model(&Asset{}).Count(&c).Error

	if countErr != nil {
		logrus.Fatal(countErr)
	}

	return c
}
