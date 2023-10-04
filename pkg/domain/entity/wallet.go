package entity

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model `json:"-"`
	ID         uint    `gorm:"primarykey"`
	UserID     int     `gorm:"userid"`
	Credited   float64 `gorm:"credit"`
	Debited    float64 `gorm:"debit"`
	Balance    float64 `gorm:"balance"`
}
