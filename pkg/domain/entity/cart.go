package entity

import "gorm.io/gorm"

type Cart struct {
	gorm.Model    `json:"-"`
	UserID        int     `json:"userid" gorm:"not null"`
	TotalProducts int     `json:"productscount"`
	TotalPrice    float64 `json:"totalprice"`
}

type CartItem struct {
	gorm.Model    `json:"-"`
	CartId        int     `json:"-"`
	UserID        int     `json:"userid" gorm:"not null"`
	ProductItemID int     `json:"productid" gorm:"not null"`
	Qty           int     `json:"qty" gorm:"not null"`
	ProductName   string  `json:"productname"`
	Price         float64 `json:"price"`
}
