package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	ID         uint   `gorm:"primarykey" bson:"_id,omitempty" json:"-"`
	FirstName  string `json:"firstname" bson:"firstname" binding:"required"`
	LastName   string `json:"lastname" bson:"lastname" binding:"required"`
	Email      string `json:"email" bson:"email" binding:"required"`
	Phone      string `json:"phone" bson:"phone" binding:"required"`
	Password   string `json:"-" bson:"password" binding:"required"`
	Wallet     int    `json:"wallet"`
	Permission bool   `gorm:"not null;default:true" json:"-"`
	Blocked    bool   `gorm:"not null;default:false" json:"-"`
}

type Address struct {
	gorm.Model `json:"-"`
	UserID     int    `json:"userid" gorm:"not null"`
	HouseName  string `json:"house_name"`
	Street     string `json:"street"`
	City       string `json:"city"`
	District   string `json:"district"`
	State      string `json:"state"`
	Pincode    string `json:"pincode"`
	Landmark   string `json:"landmark"`
}

type Login struct {
	Phone    string `json:"phone" bson:"Phone" binding:"required"`
	Password string `json:"password" bson:"Password" binding:"required"`
}

type OtpKey struct {
	gorm.Model
	Key   string `json:"key"`
	Phone string `json:"phone"`
}
