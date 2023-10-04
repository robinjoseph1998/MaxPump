package route

import (
	"MAXPUMP1/pkg/api/handlers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine, adminHandler *handlers.AdminHandler) *gin.Engine {
	//Admin Registration Routes
	r.POST("/registeradmin", adminHandler.RegisterAdmin)
	r.POST("/adminloginpassword", adminHandler.AdminLoginWithPassword)
	r.POST("/adminlogin", adminHandler.AdminLogin)
	r.POST("/adminotpvalidation", adminHandler.LoginOtpValidation)
	r.POST("/adminlogout", adminHandler.AdminLogout)

	//Users Handling routes
	r.GET("/allusers", adminHandler.ListUsers)
	r.GET("/searchuser", adminHandler.SearchUser)
	r.POST("/userblock", adminHandler.BlockUser)
	r.POST("/userunblock", adminHandler.UnBlockUser)

	//Category Handling Routes
	r.POST("/createcategory", adminHandler.CreateCategory)
	r.GET("/listcategory", adminHandler.ListCategories)
	r.PATCH("/updatecategory", adminHandler.UpdateCategory)
	r.POST("/deletecategory", adminHandler.DeleteCategory)

	//Product Handling Routes
	r.POST("/createproduct", adminHandler.AddProduct)
	r.GET("/listproducts", adminHandler.ListAllProducts)
	r.PATCH("/updateproduct", adminHandler.UpdateProduct)
	r.POST("/deleteproduct", adminHandler.DeleteProduct)
	r.GET("/loadimage", adminHandler.ShowProductPicture)

	//Coupon Handling Routes
	r.POST("/addcoupon", adminHandler.AddCoupon)
	r.GET("/showcoupons", adminHandler.ShowAllCoupons)
	r.PATCH("/editcoupon", adminHandler.EditCoupon)
	r.DELETE("/deletecoupon", adminHandler.DeleteCoupon)

	//Order Management Routes
	r.GET("/showorders", adminHandler.ViewAllOrders)
	r.POST("/changeorderstatus", adminHandler.ChangeOrderStatus)
	r.GET("/showreturnrequests", adminHandler.ShowReturnRequets)
	r.POST("/approvereturnrequest", adminHandler.ReturnRequestApproval)

	//Sales and Report Dashboard
	r.GET("/showsales", adminHandler.ViewSales)
	r.POST("/showsalesindate", adminHandler.ShowSalesOnParticularDate)
	r.POST("/salesreportpdf", adminHandler.CreateSalesReport)

	return r // Return the gin.Engine instance
}
