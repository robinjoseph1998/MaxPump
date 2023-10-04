package handlers

import (
	"MAXPUMP1/pkg/api/middleware"
	"MAXPUMP1/pkg/domain/entity"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"MAXPUMP1/pkg/utils"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type AdminHandler struct {
	AdminUsecase    use.AdminUsecaseInterface    // Adding the AdminUsecase dependency
	CategoryUsecase use.CategoryUsecaseInterface // Adding the CategoryUsecase dependency
	ProductUsecase  use.ProductUsecaseInterface  // Adding the ProductUsecase dependency
	OrderUsecase    use.OrderUsecaseInterface    //Adding the OrderUsecase dependency
}

func NewAdminHandler(adminUsecase use.AdminUsecaseInterface, categoryUsecase use.CategoryUsecaseInterface, ProductUsecase use.ProductUsecaseInterface, OrderUsecase use.OrderUsecaseInterface) *AdminHandler {
	return &AdminHandler{
		AdminUsecase:    adminUsecase,    // Initializing the AdminUsecase dependency
		CategoryUsecase: categoryUsecase, // Initializing the CategoryUsecase dependency
		ProductUsecase:  ProductUsecase,  // Initializing the ProductUsecase dependency
		OrderUsecase:    OrderUsecase,    // Initializing the OrderUsecase dependency
	}
}

// Admin Signup,Login,Logout

// @Summary New Admin Register
// @Description Registering New Admin
// @Tags Admin Registration and Login
// @Accept json
// @Produce json
// @Success 201 {object} entity.Admin "Newly Registred Admin"
// @Failure 500 "Internal Server Error"
// @Failure 400 "Bad Request"
// @Router /registeradmin [post]
func (ar AdminHandler) RegisterAdmin(c *gin.Context) {
	var admin entity.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAdmin, err := ar.AdminUsecase.ExecuteAdminCreate(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newAdmin)
}

// @Summary Admin Login With Password
// @Description Admin Login using password and phone number
// @Tags Admin Registration and Login
// @Accept json
// @Produce json
// @Param User body utils.AdminLogin true "Phone and Password"
// @Success 200 "message":"Admin Logged in Successfully"
// @Failure 400 "Bad Request"
// @Failure 400 "Bad Request"
// @Router /adminloginpassword [post]
func (ar *AdminHandler) AdminLoginWithPassword(c *gin.Context) {
	var loginrequest utils.AdminLogin
	if err := c.ShouldBindJSON(&loginrequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	phone := loginrequest.Phone
	password := loginrequest.Password
	adminId, err := ar.AdminUsecase.ExecuteLoginWithPassword(phone, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		middleware.GenToken(uint(adminId), phone, c)
		c.JSON(http.StatusOK, gin.H{"message": "Admin logged in succesfully and cookie stored"})
	}
}

// @Summary Admin Login
// @Description Admin Login Otp
// @Tags Admin Registration and Login
// @Accept json
// @Produce json
// @Param User body utils.LoginPayload true "Phone"
// @Success 200 "Otp send to"":Payload.Phone
// @Failure 401 "Unauthorized"
// @Failure 400 "Bad Request"
// @Router /adminlogin [post]
func (ar *AdminHandler) AdminLogin(c *gin.Context) {
	var PayLoad utils.LoginPayload
	if err := c.ShouldBindJSON(&PayLoad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ar.AdminUsecase.ExecuteAdminLogin(PayLoad.Phone); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Otp send succesfully to": PayLoad.Phone})
	}
}

// @Summary Login Otp Validation
// @Description Validating Entered Otp By User
// @Tags Admin Registration and Login
// @Accept json
// @Produce json
// @Param User body utils.OtpKey true "Phone,Otp and Resend"
// @Success 200 "message":"Admin Logged in Successfully and cookie stored"
// @Failure 401 "Unauthorized"
// @Success 200 "message":"OTP resend successful"
// @Failure 400 "Bad Request"
// @Router /adminotpvalidation [post]
func (ar *AdminHandler) LoginOtpValidation(c *gin.Context) {
	var otpvalidation utils.OtpKey
	if err := c.ShouldBindJSON(&otpvalidation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if otpvalidation.Resend == "resend" {
		err := ar.AdminUsecase.ExecuteAdminLogin(otpvalidation.Phone)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "OTP resend successful"})
		}
	} else {
		admin, err := ar.AdminUsecase.ExecuteOtpValidation(otpvalidation.Phone, otpvalidation.OTP)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		middleware.GenToken(admin.ID, admin.Phone, c)
		c.JSON(http.StatusOK, gin.H{"message": "Admin logged in successfully and cookie stored"})
	}
}

// @Summary Admin Logout
// @Description Logout By Admin
// @Tags Admin Registration and Login
// @Accept json
// @Produce json
// @Success 200 "message":"logged out successfully"
// @Failure 400 "Bad Request"
// @Router /adminlogout [post]
func (ar *AdminHandler) AdminLogout(c *gin.Context) {
	err := middleware.DeleteCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user cookie deletion failed"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
	}
}

// Admin Side User Management

// @Summary List Users
// @Description List All Users in Admin Side
// @Tags Admin Side User Management
// @Accept json
// @Produce json
// @Success 200 "users":users
// @Failure 500 "Internal Server Error"
// @Router /allusers [get]
func (ar *AdminHandler) ListUsers(c *gin.Context) {
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number required"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page Number"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	TotalUsers, err := ar.AdminUsecase.ExecuteTotalOfUsers()
	totalPages := (TotalUsers + itemsPerPage - 1) / itemsPerPage
	paginatedUsers, err := ar.AdminUsecase.ExecuteAllUsersPaginated(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{
		"PaginatedUsers": paginatedUsers,
		"TotalUsers":     TotalUsers,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Users": response})
}

// @Summary Search User
// @Description Search User By id and firstname
// @Tags Admin Side User Management
// @Accept json
// @Produce json
// @Params User body utils.SearchUser true "ID and FirstName"
// @Sucess 200 "user":user
// @Failure 500 "Internal Server Error"
// @Failure 400 "Bad Request"
// @Router /searchuser [get]
func (ar *AdminHandler) SearchUser(c *gin.Context) {
	var Search utils.SearchUser
	if err := c.ShouldBindJSON(&Search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := ar.AdminUsecase.ExecuteSearch(Search.ID, Search.FirstName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// @Summary Block User
// @Description Block User By Admin
// @Tags Admin Side User Management
// @Accept json
// @Produce json
// @Params User body utils.BlockUser true "ID"
// @Success 200 "message":"User Blocked Successfully"
// @Failure 500 "Internal Server Error"
// @Failure 400 "Bad Request"
// @Router /userblock [post]
func (ar *AdminHandler) BlockUser(c *gin.Context) {
	var blocker utils.BlockUser
	if err := c.ShouldBindJSON(&blocker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ar.AdminUsecase.BlockUserByID(blocker.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// @Summary Unblock User
// @Description Unblock user by admin
// @Tags Admin Side User Management
// @Accept json
// @Produce json
// @Params User body utils.BlockUser true "ID"
// @Success 200 "message":"User Unblocked Successfully"
// @Failure 500 "Internal Server Error"
// @Failure 400 "Bad Request"
// @Router /userunblock [post]
func (ar *AdminHandler) UnBlockUser(c *gin.Context) {
	var Unblocker utils.BlockUser
	if err := c.ShouldBindJSON(&Unblocker); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ar.AdminUsecase.UnBlockUserByID(Unblocker.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Unblocked successfully"})
}

//Admin Side Category Management

// @Summary Create Category
// @Description Category Creation by Admin
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Params User body entity.Category true "ID,Name,Description and view"
// @Success 201 "Category":category
// @Failure 500 "Internal Server Error"
// @Failure 500 "Internal Server Error"
// @Router /createcategory [post]
func (ar *AdminHandler) CreateCategory(c *gin.Context) {
	var cat entity.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	category, err := ar.CategoryUsecase.CategoryCreater(cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Category": category})
}

// @Summary List Categories
// @Description Admin Listing Categories
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Success 200 "Categories":categories
// @Failure 500 "Internal Server Error"
// @Router /listcategory [get]
func (ar *AdminHandler) ListCategories(c *gin.Context) {
	categories, CategoryIdAndProductsCount, CategoryIdAndBrandsCount, err := ar.CategoryUsecase.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var response []gin.H
	for _, category := range categories {
		categoryResponse := gin.H{
			"Category ID":   category.ID,
			"Category Name": category.Name,
			"About":         category.Description,
		}
		for _, CategoryIdAndCount := range CategoryIdAndProductsCount {
			if count, exists := CategoryIdAndCount[int(category.ID)]; exists {
				categoryResponse["Total Products Items"] = count
			}
		}
		for _, CategoryIdAndBrandCount := range CategoryIdAndBrandsCount {
			if brandscount, exists := CategoryIdAndBrandCount[int(category.ID)]; exists {
				categoryResponse["Total Brands"] = brandscount
			}
		}
		response = append(response, categoryResponse)
	}
	c.JSON(http.StatusOK, gin.H{"Categories": response})
}

// @Summary Update Category
// @Description Admin Update Category
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Params User body utils.EditCategory true "ID,Name and Description"
// @Success 200 "Updated Category":category
// @Failure 500 "Internal Server Error"
// @Failure 404 "error":"category not found"
// @Failure 500 "Internal Server Error"
// @Router /updatecategory [patch]
func (ar *AdminHandler) UpdateCategory(c *gin.Context) {
	var EditCategory utils.EditCategory
	if err := c.ShouldBindJSON(&EditCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	category, err := ar.CategoryUsecase.GetCategoryByID(EditCategory.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category from database"})
		return
	}
	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	category.Name = EditCategory.Name
	category.Description = EditCategory.Description
	category, err = ar.CategoryUsecase.CategoryUpdate(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Updated Category": category})
}

// @Summary Delete Category
// @Description Admin Delete Category
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Params User body utils.DeleteCategory true "ID"
// @Success 200 "message":"Category deleted successfully"
// @Failure 500 "Internal Server Error"
// @Failure 500 "Internal Server Error"
// @Failure 400 "Bad Request"
// @Router /deletecategory [post]
func (ar *AdminHandler) DeleteCategory(c *gin.Context) {
	var deletecategory utils.DeleteCategory
	if err := c.ShouldBindJSON(&deletecategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	category, err := ar.CategoryUsecase.GetCategoryByID(deletecategory.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = ar.CategoryUsecase.DeleteCategory(category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// Admin Side Product Management

// @Summary Add Product
// @Description Admin Product Adding
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Params User body entity.product true "ID,Name,Description,Price,Quantity and CategoryID"
// @Success 201 "responseData"
// @Failure 500 "Internal Server Error"
// @Failure 500 "Internal Server Error"
// @Router /createproduct [post]
func (ar *AdminHandler) AddProduct(c *gin.Context) {
	BrandName := c.PostForm("brandname")
	Description := c.PostForm("description")
	Item := c.PostForm("item")
	PriceStr := c.PostForm("price")
	QuantityStr := c.PostForm("qty")
	CategoryIDStr := c.PostForm("category_id")
	TypeConvertedPrice, err := strconv.ParseFloat(PriceStr, 64)
	if err != nil || TypeConvertedPrice <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}
	TypeConvertedQuantity, err := strconv.Atoi(QuantityStr)
	if err != nil || TypeConvertedQuantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}
	TypeConvertedCategoryID, err := strconv.Atoi(CategoryIDStr)
	if err != nil || TypeConvertedCategoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid image"})
		return
	}
	imageFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer imageFile.Close()
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO()) //Loading the default AWS sdk configuration
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	client := s3.NewFromConfig(cfg)                                    //creating an s3 client
	uploader := manager.NewUploader(client)                            // setup a uploader
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{ //uploads the image file to an S3 bucket named "maxpumpbucket"
		Bucket: aws.String("maxpumpbucket"), //bucket name
		Key:    aws.String(file.Filename),   //object key sets as file name
		Body:   imageFile,                   //object content sets as img file
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := entity.Product{
		Brand_Name:  BrandName,
		Description: Description,
		Item:        Item,
		Price:       TypeConvertedPrice,
		Quantity:    TypeConvertedQuantity,
		CategoryID:  TypeConvertedCategoryID,
		ImageURL:    result.Location,
	}
	AddedProduct, category, err := ar.ProductUsecase.ProductCreater(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Product":  AddedProduct,
		"Category": category,
	})
}

// @Summary List Products
// @Description Admin List Products
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Success 200 "Products":AllProducts
// @Failure 500 "Internal Server Error"
// @Router /listproducts [get]
func (ah *AdminHandler) ListAllProducts(c *gin.Context) {
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number required"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page Number"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	PaginatedProducts, err := ah.ProductUsecase.ExecutePaginatedProducts(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(PaginatedProducts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no products to show"})
		return
	}
	if PaginatedProducts == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no products to show"})
		return
	}
	TotalProducts, err := ah.ProductUsecase.ExecuteTotalProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get the total number of products"})
		return
	}
	totalPages := (TotalProducts + itemsPerPage - 1) / itemsPerPage
	response := gin.H{
		"PaginatedProducts": PaginatedProducts,
		"TotalProducts":     TotalProducts,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Products": response})
}

// @Summary Update Product
// @Description Product Updation by Admin
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Params User body utils.EditProduct true "ID,Name,Description and Price"
// @Success 200 "Product":Resultproduct
// @Failure 500 "Internal Server Error"
// @Failure 404 "error":"product not found"
// @Failure 500 "error":"Failed to fetch product from database"
// @Failure 400 "Bad Request"
// @Router /updateproduct [post]
func (ar *AdminHandler) UpdateProduct(c *gin.Context) {
	var EditProduct *utils.EditProductRequest
	if err := c.ShouldBindJSON(&EditProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	product, category, err := ar.ProductUsecase.ProductUpdate(EditProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Resultproduct := gin.H{
		"Product":  product,
		"Category": category,
	}
	c.JSON(http.StatusOK, gin.H{"Product": Resultproduct})
}

// @Summary Delete Product
// @Description Admin Product Delete
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Params User body utils.DeletingProduct true "ID"
// @success 200 "message":"Product Deleted successfully"
// @Failure 500 "Internal Server Error"
// @Failure 500 "Internal Server Error"
// @Router /deleteproduct [post]
func (ar *AdminHandler) DeleteProduct(c *gin.Context) {
	var product utils.DeletingProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err := ar.ProductUsecase.DeleteProduct(product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Product Deleted Successfully"})
}

//Coupon Management

func (ar *AdminHandler) AddCoupon(c *gin.Context) {
	var coupon *entity.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	CreatedCoupon, err := ar.AdminUsecase.CreateCoupon(coupon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Message": "Coupon Created Successfully",
		"Coupon":  CreatedCoupon})
}

func (ar *AdminHandler) ShowAllCoupons(c *gin.Context) {
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number required"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page Number"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	TotalCoupons, err := ar.AdminUsecase.ExecuteTotalOfCoupons()
	if TotalCoupons == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no coupons available"})
		return
	}
	totalPages := (TotalCoupons + itemsPerPage - 1) / itemsPerPage
	PaginatedCoupons, err := ar.AdminUsecase.ExecutePaginatedCoupons(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't show coupons"})
		return
	}
	if PaginatedCoupons == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no coupons to show"})
		return
	}
	response := gin.H{
		"PaginatedCoupons": PaginatedCoupons,
		"TotalCoupons":     TotalCoupons,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Coupons": response})
}

func (ar *AdminHandler) EditCoupon(c *gin.Context) {
	var request *utils.EditCouponRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	EditedCoupon, err := ar.AdminUsecase.ExecuteEditCoupon(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Coupon": EditedCoupon})
}

func (ar *AdminHandler) DeleteCoupon(c *gin.Context) {
	var DeletingCouponCode utils.CouponCode
	if err := c.ShouldBind(&DeletingCouponCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	code := DeletingCouponCode.Code
	if err := ar.AdminUsecase.ExecuteDeleteCoupon(code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "coupon deleted successfully"})
}

//Order Management

func (ar *AdminHandler) ViewAllOrders(c *gin.Context) {
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number required"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page Number"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	TotalOrders, err := ar.OrderUsecase.ExecuteTotalOfAllOrders()
	if TotalOrders == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "currently no sales"})
		return
	}
	totalPages := (TotalOrders + itemsPerPage - 1) / itemsPerPage
	PaginatedOrders, err := ar.AdminUsecase.ExecuteAllOrdersPaginated(offset, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't fetch orders"})
		return
	}
	if PaginatedOrders == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no orders to show"})
		return
	}
	response := gin.H{
		"PaginatedOrders": PaginatedOrders,
		"TotalOrders":     TotalOrders,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Coupons": response})
}

func (ar *AdminHandler) ChangeOrderStatus(c *gin.Context) {
	var ApplyingStatus *utils.OrderStatus
	if err := c.ShouldBindJSON(&ApplyingStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch ApplyingStatus.Status {
	case "Delivered", "Shipped", "Processing", "Cancelled", "Returned":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status option"})
		return
	}
	StatusUpdatedOrder, err := ar.AdminUsecase.ChangeStatus(ApplyingStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't change status"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "status updated successfully",
		"Order":   StatusUpdatedOrder})
}

func (ar *AdminHandler) ShowReturnRequets(c *gin.Context) {
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number required"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page Number"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	TotalRequests, err := ar.OrderUsecase.ExecuteTotalOfReturnRequests()
	if TotalRequests == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "currently no Return Requests"})
		return
	}
	totalPages := (TotalRequests + itemsPerPage - 1) / itemsPerPage
	ReturnRequestsPaginated, err := ar.OrderUsecase.ExecuteShowReturnRequestsPaginated(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ReturnRequestsPaginated == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no requests to show"})
		return
	}
	response := gin.H{
		"PaginatedreturnRequests": ReturnRequestsPaginated,
		"TotalRequests":           TotalRequests,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Return Requests": response})
}

func (ar *AdminHandler) ReturnRequestApproval(c *gin.Context) {
	var Request *utils.ReturnProductRequest
	if err := c.ShouldBind(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind approval request"})
		return
	}
	Order, err := ar.AdminUsecase.ExecuteReturnRequestApproval(Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Returned order": Order})
}

// Dashboard Management
func (ar *AdminHandler) ViewSales(c *gin.Context) {
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number required"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page Number"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	TotalSales, err := ar.OrderUsecase.ExecuteTotalOfSales()
	if TotalSales == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "currently no sales"})
		return
	}
	totalPages := (TotalSales + itemsPerPage - 1) / itemsPerPage
	PaginatedSales, err := ar.OrderUsecase.ExecuteSuccessfullOrdersPaginated(offset, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't view sales"})
		return
	}
	if PaginatedSales == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no sales to show"})
		return
	}
	response := gin.H{
		"PaginatedSales": PaginatedSales,
		"TotalSales":     TotalSales,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Sales": response})
}

func (ar *AdminHandler) ShowSalesOnParticularDate(c *gin.Context) {
	var request utils.MyDate
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind the request"})
		return
	}
	dateStr := strings.TrimSpace(request.Date)
	layout := "2006-01-02"
	ParsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page number required"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Page Number"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	TotalSalesByDate, err := ar.OrderUsecase.ExecuteTotalOfSalesByDate(ParsedTime)
	if TotalSalesByDate == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no sales on this date"})
		return
	}
	totalPages := (TotalSalesByDate + itemsPerPage - 1) / itemsPerPage
	PaginatedSales, err := ar.OrderUsecase.ExecuteSalesInParticularDatePaginated(ParsedTime, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't fetch sales in this date"})
		return
	}
	if PaginatedSales == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no sales to show"})
		return
	}
	response := gin.H{
		"Date":           ParsedTime,
		"PaginatedSales": PaginatedSales,
		"TotalSales":     TotalSalesByDate,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Sales": response})
}

func (ar *AdminHandler) CreateSalesReport(c *gin.Context) {
	var request utils.MyDate
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind the date"})
		return
	}
	dateStr := strings.TrimSpace(request.Date)

	err := ar.AdminUsecase.ExecuteCreateSalesReportByDate(dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't generate report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sales Report generated successfully"})

}

func (ar *AdminHandler) ShowProductPicture(c *gin.Context) {
	imageUrl := c.DefaultQuery("url", "")
	c.HTML(http.StatusOK, "webpage.html", gin.H{
		"ImageUrl": imageUrl, // Passing the image URL to the template
	})
}
