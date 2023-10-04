package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model  `json:"-"`
	ID          uint    `gorm:"primarykey"`
	Brand_Name  string  `json:"brandname" gorm:"not null"`
	Description string  `json:"description"`
	Item        string  `json:"item"`
	Price       float64 `json:"price" gorm:"not null"`
	Quantity    int     `json:"qty" gorm:"not null"`
	CategoryID  int     `json:"category_id"`
	ImageURL    string  `json:"imageurl"`
}
