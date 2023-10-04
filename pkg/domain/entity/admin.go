package entity

import "gorm.io/gorm"

type Admin struct {
	gorm.Model `json:"-"`
	AdminName  string `json:"adminname" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Active     bool   `gorm:"not null;default true"`
}

type AdminOtpKey struct {
	gorm.Model
	Key   string `json:"key"`
	Phone string `json:"phone"`
}
