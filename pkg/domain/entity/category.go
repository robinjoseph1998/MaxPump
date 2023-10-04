package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model  `json:"-"`
	ID          uint   `gorm:"primarykey" bson:"_id,omitempty" json:"-"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	// view        bool   `gorm:"not null;default true"`
}
