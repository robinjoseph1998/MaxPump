package route

import (
	"MAXPUMP1/pkg/api/handlers"
	"MAXPUMP1/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) *gin.Engine {
	//User Signup and Login Routes
	r.GET("/welcome", userHandler.WelcomeMessage)
	r.POST("/signup", userHandler.Signup)
	r.POST("/signupotp", userHandler.SignupWithOtp)
	r.POST("/signupotpvalidation", userHandler.SignupOtpValidation)
	r.POST("/loginwithpassord", userHandler.LoginWithPassword)
	r.POST("/logout", userHandler.UserLogout)

	//User Categories and Product routes
	r.GET("/categoriesview", userHandler.CategoriesView)
	r.GET("/productsbycategory", userHandler.ProductsByCategory)
	r.GET("/allproducts", userHandler.ListAllProducts)
	r.GET("/productdetail", userHandler.ProductDetail)

	//User Cart routes
	r.POST("/addtocart", middleware.ValidateCookie, userHandler.AddToCart)
	r.GET("/listcart", middleware.ValidateCookie, userHandler.ListCart)
	r.DELETE("/deletecartproduct", middleware.ValidateCookie, userHandler.DeleteProductFromCart)

	//User Address & Profile Handling Routes
	r.POST("/addaddress", middleware.ValidateCookie, userHandler.AddAddress)
	r.GET("/viewaddresses", middleware.ValidateCookie, userHandler.ViewAddresses)
	r.PATCH("/editaddress", middleware.ValidateCookie, userHandler.EditAddress)
	r.PUT("/updateaddress", middleware.ValidateCookie, userHandler.UpdateAddress)
	r.GET("/showprofile", middleware.ValidateCookie, userHandler.ShowProfile)

	//User Checkout Order Routes
	r.POST("/placeorder", middleware.ValidateCookie, userHandler.PlaceOrder)
	r.POST("/paymentverification", middleware.ValidateCookie, userHandler.VerifyMyPayment)
	r.GET("/vieworder", middleware.ValidateCookie, userHandler.ViewOrder)
	r.POST("/cancelorder", middleware.ValidateCookie, userHandler.CancellMyOrder)
	r.PATCH("/returnproduct", middleware.ValidateCookie, userHandler.Returnproduct)
	r.POST("/generateinovicepdf", middleware.ValidateCookie, userHandler.CreateInvoice)

	//User Coupon Management
	r.GET("/availablecoupons", middleware.ValidateCookie, userHandler.ShowAvailableCoupons)
	r.POST("/applycoupon", middleware.ValidateCookie, userHandler.ApplyCoupon)

	//User Product And Category Filtering
	r.POST("/getcategory", middleware.ValidateCookie, userHandler.FilterCategory)
	r.POST("/getproductsbybrand", middleware.ValidateCookie, userHandler.FilterProductByBrand)
	r.POST("/getproductsbyitem", middleware.ValidateCookie, userHandler.FilterProductByItem)

	//User Wallet
	r.GET("/showmywallet", middleware.ValidateCookie, userHandler.ShowWallet)

	return r
}
