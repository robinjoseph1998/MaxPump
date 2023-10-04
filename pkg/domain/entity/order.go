package entity

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model        `json:"-"`
	ID                int     `gorm:"primarykey" json:"id"`
	UserID            int     `json:"userid"`
	AddressId         int     `json:"adressid"`
	TotalPrice        float64 `json:"totalprice"`
	PaymentMethod     string  `json:"paymentmethod"`
	PaymentStatus     string  `json:"paymentstatus"`
	OrderStatus       string  `json:"Orderstatus"`
	Date_Of_Delivered string  `json:"dod"`
}

type Ordered_Item struct {
	gorm.Model `json:"-"`
	Id         int     `gorm:"primarykey" json:"id"`
	ProductID  int     `json:"productid"`
	UserID     int     `json:"userid"`
	OrderID    int     `json:"orderid"`
	Item       string  `json:"item"`
	Quantity   int     `json:"qty"`
	Price      float64 `json:"price"`
	Status     string  `json:"Orderstatus"`
}

type ReturnedOrder struct {
	gorm.Model        `json:"-"`
	ID                uint    `gorm:"primarykey"`
	Price             float64 `gorm:"price"`
	Item              string  `gorm:"item"`
	ProductID         int     `gorm:"productid"`
	Orderid           int     `gorm:"orderid"`
	UserID            int     `gorm:"userid"`
	Status            string  `gorm:"status"`
	Date_Of_Delivered string  `json:"dod"`
}
