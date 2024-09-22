package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Helps to fetch User Id from Request Context
func GetUserIDFromContext(c *gin.Context) (int, error) {
	userIdStr := c.GetString("userID")
	fmt.Println("userID", userIdStr)
	userID, err := strconv.Atoi(userIdStr)
	return userID, err
}

type EditAddressRequest struct {
	HouseName *string `json:"house_name"`
	Street    *string `json:"street"`
	City      *string `json:"city"`
	District  *string `json:"district"`
	State     *string `json:"state"`
	Pincode   *string `json:"pincode"`
	Landmark  *string `json:"landmark"`
}

type OrderRequest struct {
	UserID        int    `json:"userid"`
	AddressId     int    `json:"addressid"`
	PaymentMethod string `json:"paymentmethod"`
}

type OtpValidation struct {
	Key string `json:"key"`
	Otp string `json:"otp"`
}

type UserLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type CatAndPro struct {
	Name string `json:"name" binding:"required"`
}

type Product struct {
	ID uint `json:"id" binding:"required"`
}

type CartRequest struct {
	ProductItemID int `json:"productid" gorm:"not null"`
	Qty           int `json:"qty" gorm:"not null"`
}

type ProductDrop struct {
	ProductID int `json:"productid" gorm:"not null"`
}

type OrderID struct {
	ID int `json:"id"`
}

type AdminLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Phone string `json:"phone"`
}

type OtpKey struct {
	Phone  string `json:"phone"`
	OTP    string `json:"otp"`
	Resend string `json:"resend"`
}

type SearchUser struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
}

type BlockUser struct {
	ID uint `json:"id" binding:"required"`
}

type EditCategory struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type DeleteCategory struct {
	ID uint `json:"id" binding:"required"`
}

type EditProduct struct {
	ID          uint    `json:"id" binding:"required"`
	Brand_Name  string  `json:"brandname"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"qty"`
}

type DeletingProduct struct {
	ID uint `json:"id" binding:"required"`
}

type CouponCode struct {
	Code int `json:"code" binding:"required"`
}

type OrderStatus struct {
	ID     int    `json:"orderid" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type CategoryFilterRequest struct {
	ID int `json:"id"`
}

type EditProductRequest struct {
	ID          uint     `json:"id"`
	Brand_Name  *string  `json:"brandname"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Quantity    *int     `json:"qty"`
	Category    *int     `json:"category"`
}

type EditCouponRequest struct {
	Code             int        `json:"code"`
	Type             *string    `json:"type"`
	Amount           *float64   `json:"amount"`
	Threshold_Amount *float64   `json:"threshold"`
	Expiration       *time.Time `json:"expiration"`
	UsageLimit       *int       `json:"usage_limit"`
}

type ReturnProductRequest struct {
	OrderID   int `json:"orderid"`
	ProductID int `json:"productid"`
}

type PaymentVerificationRequest struct {
	PaymentID string `json:"paymentid"`
	OrderID   string `json:"orderid"`
	Signature string `json:"sign"`
}

type MyDate struct {
	Date string `json:"date"`
}

type InvoiceModel struct {
	InvoiceNumber  string
	Date           string
	BillingName    string
	BillingAddress string
	District       string
	Pincode        string
	Landmark       string
	TotalPrice     float64
}

type SalesReportModel struct {
	Date             time.Time
	TotalSales       int
	TotalAmount      float64
	ProductWiseSales map[string]int
	PricePerQuantity map[string]float64
}

func GenerateUniqueKey() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func ImageLoader(c *gin.Context) {
	htmlFilePath := "assets/index.html"
	c.File(htmlFilePath)
}
