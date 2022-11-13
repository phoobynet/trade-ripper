package calendars

import (
	"gorm.io/gorm"
	"time"
)

type Calendar struct {
	gorm.Model   `json:"-"`
	Date         string    `json:"date" gorm:"uniqueIndex"`
	Open         time.Time `json:"open"`
	SessionOpen  time.Time `json:"sessionOpen"`
	Close        time.Time `json:"close"`
	SessionClose time.Time `json:"sessionClose"`
}
