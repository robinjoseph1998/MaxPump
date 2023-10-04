package db

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/model"

	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB() *gorm.DB {

	dsn := "host=localhost user=postgres password=robin123 dbname=maxpump port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting database")
		return nil
	}

	DB = db

	DB.AutoMigrate(&entity.User{})
	DB.AutoMigrate(&entity.OtpKey{})
	DB.AutoMigrate(&model.Signup{})
	DB.AutoMigrate(&entity.Admin{})
	DB.AutoMigrate(&entity.AdminOtpKey{})
	DB.AutoMigrate(&entity.Category{})
	DB.AutoMigrate(&entity.Product{})
	DB.AutoMigrate(&entity.Cart{})
	DB.AutoMigrate(&entity.CartItem{})
	DB.AutoMigrate(&entity.Address{})
	DB.AutoMigrate(&entity.Order{})
	DB.AutoMigrate(&entity.Coupon{})
	DB.AutoMigrate(&entity.ReturnedOrder{})
	DB.AutoMigrate(&entity.Ordered_Item{})
	DB.AutoMigrate(&entity.Wallet{})

	return db
}
