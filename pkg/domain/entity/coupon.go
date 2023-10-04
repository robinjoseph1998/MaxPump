package entity

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model       `json:"-"`
	Id               int       `json:"-" gorm:"primarykey"`
	Code             int       `json:"code"`
	Type             string    `json:"type"`
	Amount           float64   `json:"amount"`
	Threshold_Amount float64   `json:"threshold"`
	Expiration       time.Time `json:"expiration"`
	UsageLimit       int       `json:"usage_limit"`
	UsedCount        int       `json:"used_count"`
}
