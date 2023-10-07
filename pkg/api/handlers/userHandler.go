package handlers

import (
	"MAXPUMP1/pkg/api/middleware"
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/model"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"MAXPUMP1/pkg/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserHandler struct {
	UserUsecase     use.UserUsecaseInterface     // Adding the UserUsecase dependency
	CategoryUsecase use.CategoryUsecaseInterface // Adding the CategoryUsecase dependency
	ProductUsecase  use.ProductUsecaseInterface  // Adding the ProductUsecase dependency
	CartUsecase     use.CartUsecaseInterface     // Adding the CartUsecase dependency
	OrderUsecase    use.OrderUsecaseInterface    // Adding the OrderUsecase dependency
}

func NewUserHandler(UserUsecase use.UserUsecaseInterface, categoryUsecase use.CategoryUsecaseInterface, ProductUsecase use.ProductUsecaseInterface, CartUsecase use.CartUsecaseInterface, OrderUsecase use.OrderUsecaseInterface) *UserHandler {
	return &UserHandler{UserUsecase: UserUsecase, // Initializing the AdminUsecase dependency
		CategoryUsecase: categoryUsecase, // Initializing the CategoryUsecase dependency
		ProductUsecase:  ProductUsecase,  // Initializing the ProductUsecase dependency
		CartUsecase:     CartUsecase,     // Initializing the CartUsecase dependency
		OrderUsecase:    OrderUsecase,    // Initializing the OrderUsecase dependency
	}
}

// Signup handles user signup.
//
//	@Summary		User Signup
//	@Description	Register a new user.
//	@Tags			User Authentication
//	@Accept			json
//	@Produce		json
//	@Param			userInput	body		model.Signup	true	"User Signup Input"
//	@Success		201			{object}	entity.User		"Newly registered user"
//	@Failure		400			"Bad Request"
//	@Failure		500			"Internal Server Error"
//	@Router			/signup [post]
func (uh *UserHandler) Signup(c *gin.Context) {
	var userInput model.Signup
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	copier.Copy(&user, &userInput)
	newUser, err := uh.UserUsecase.ExecuteSignup(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

// @Summary		signup with otp
// @Description	Adding new user to the database
// @Tags			User Authentication
// @Accept			json
// @Produce		json
// @Param			User	body		model.Signup	true	"User Data"
// @Success		200		{object}	entity.User
// @Failure		400		"Bad Request"
// @Failure		500		"Internal Server Error"
// @Router			/signupotp [post]
func (uh *UserHandler) SignupWithOtp(c *gin.Context) {
	var user model.Signup
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := uh.UserUsecase.ExecuteSignupWithOtp(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Otp send succesfuly to": user.Phone, "Key": key})
	}
}

// @Summary		Signup OTP Validation
// @Description	Validating entered OTP by user
// @Tags			User Authentication
// @Accept			json
// @Produce		json
// @Param			User	body		utils.OtpValidation	true	"Key and OTP"
// @Success		200		"message":	"user signup successful"
// @Failure		400		"error":	"Bad Request"
// @Failure		401		"error":	"Unauthorized"
// @Router			/signupotpvalidation [post]
func (uh *UserHandler) SignupOtpValidation(c *gin.Context) {
	var request utils.OtpValidation
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := request.Key
	otp := request.Otp
	err := uh.UserUsecase.ExecuteSignupOtpValidation(key, otp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user signup succesfull"})
	}
}

// @Summary		Login With Password
// @Description	Validating entered Password by user
// @Tags			User Authentication
// @Accept			json
// @Produce		json
// @Param			User	body		utils.UserLogin	true	"Phone and Password"
// @Success		200		"message":	"user logged in succesfully and cookie stored"
// @Failure		400		"error":	"Bad Request"
// @Failure		404		"error":	"Status Not Found"
// @Router			/loginwithpassord [post]
func (uh *UserHandler) LoginWithPassword(c *gin.Context) {
	var loginrequest utils.UserLogin
	if err := c.ShouldBindJSON(&loginrequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	phone := loginrequest.Phone
	password := loginrequest.Password
	userID, err := uh.UserUsecase.ExecuteLoginWithPassword(phone, password)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else {
		middleware.GenToken(userID, phone, c)
		c.JSON(http.StatusOK, gin.H{"message": "user logged in succesfully and cookie stored"})
	}
}

// @Summary Logout
// @Description User Logout
// @Tags User Authentication
// @Accept json
// @Produce json
// @Success 200 "message":"logged out successfully"
// @Failure 400 "error":"user cookie deletion failed"
// @Router /logout [post]
func (uh *UserHandler) UserLogout(c *gin.Context) {
	err := middleware.DeleteCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user cookie deletion failed"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
	}
}

// User Category Handlers

// @Summary View Categories
// @Description List All Categories
// @Tags User Category and Product Management
// @Accept json
// @Produce json
// @Success 200 "Categories":AllCategories
// @Failure 400 "error":"Internal Server Error"
// @Router /categoriesview [get]
func (uh *UserHandler) CategoriesView(c *gin.Context) {
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
	PaginatedCategories, err := uh.CategoryUsecase.ExecutePaginatedCategories(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	TotalCategories, err := uh.CategoryUsecase.ExecuteTotalCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get the total number of categories"})
		return
	}
	totalPages := (TotalCategories + itemsPerPage - 1) / itemsPerPage

	response := gin.H{
		"Categories": PaginatedCategories,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Categories": response})
}

// @Summary Category with Products
// @Description List Category With Products
// @Tags User Category and Product Management
// @Params Category body utils.CatAndPro true Name
// @Accept json
// @Produce json
// @Success 200 "category and products":CategoryWithProducts
// @Failure 500 "error":"Internal Server Error"
// @Router /productsbycategory [get]
func (uh *UserHandler) ProductsByCategory(c *gin.Context) {
	var categoryProducts utils.CatAndPro
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
	if err := c.ShouldBindJSON(&categoryProducts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	PaginatedProducts, category, err := uh.CategoryUsecase.CategoryWithPaginatedProducts(categoryProducts.Name, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	TotalCategories, err := uh.CategoryUsecase.ExecuteTotalProductsInTheParticularCategory(int(category.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get the total number of categories"})
		return
	}
	totalPages := (TotalCategories + itemsPerPage - 1) / itemsPerPage
	response := gin.H{
		"Products": PaginatedProducts,
		"Category": category,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Categories": response})
}

// User Product Handlers

// @Summary Listing Products
// @Description Listing all products
// @Tags User Category and Product Management
// @Accept json
// @Produce json
// @Success 200 "Products":AllProducts
// @Failure 500 "Internal Server Error"
// @Router /allproducts [get]
func (uh *UserHandler) ListAllProducts(c *gin.Context) {
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
	PaginatedProducts, err := uh.ProductUsecase.ExecutePaginatedProducts(offset, limit)
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
	TotalProducts, err := uh.ProductUsecase.ExecuteTotalProducts()
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

// @Summary Showing Product Details
// @Description Showing Product Detaily
// @Tags User Category and Product Management
// @Accept json
// @Produce json
// @Params Product utils.Product true ID
// @Success 200 "product":ProductDetailedView
// @Failure 500 "Internal Server Error"
// @Failure 400 "Bad Request"
// @Router /productdetail [get]
func (uh *UserHandler) ProductDetail(c *gin.Context) {
	var product utils.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ProductDetail, ProductCategory, err := uh.ProductUsecase.GetProductByID(product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product from database"})
		return
	}
	if ProductDetail == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	ProductDetailedView := gin.H{
		"Product":  ProductDetail,
		"Category": ProductCategory,
	}
	c.JSON(http.StatusOK, gin.H{"Product": ProductDetailedView})
}

//User Cart Handlers

// @Summary Adding Product To Cart
// @Description Adding Product To User Cart
// @Tags User Cart Management
// @Accept json
// @Produce json
// @Params Cart utils.CartRequest true ProductItemID and Qty
// @Success 200 "message":cart and cart Item
// @Failure 500 "error":"Internal Server Error"
// Failure  404 "message":"Status Not Found"
// Failure  500 "error":"Internal Server Error"
// Failure  400 "error":"Bad Request"
// @Router /addtocart [post]
func (uh *UserHandler) AddToCart(c *gin.Context) {
	var cart utils.CartRequest
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := uh.CartUsecase.GetProductByID(uint(cart.ProductItemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to Find The Product"})
		return
	}
	Qty, err := uh.CartUsecase.CheckQuantity(int(product.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check quantity from database"})
		return
	}
	if Qty == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Sorry This Product Is Out Of Stock"})
		return
	}
	if Qty < cart.Qty {
		c.JSON(http.StatusNotFound, gin.H{"message": "the quantinty you entered is exceeds the limit of our stock"})
		return
	}
	userID, _ := utils.GetUserIDFromContext(c)
	Cart, cartItem, err := uh.CartUsecase.AddToCart(userID, int(product.ID), cart.Qty)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Cart":       Cart,
		"Added Item": cartItem})
}

// @Summary List User Cart
// @Description Listing User's Cart
// @Tags User Cart Management
// @Accept json
// @Produce json
// @Params Cart utils.GetUserIDFromContext() true userID
// @Success 200 "message":"Cart and CartItems"
// @Failure 500 "error":"Internal Server Error"
// @Failure 500 "error":"Internal Server Error"
// @Router /listcart [get]
func (uh *UserHandler) ListCart(c *gin.Context) {
	userID, _ := utils.GetUserIDFromContext(c)
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
	Cart, err := uh.CartUsecase.ExecuteCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch Cart"})
		return
	}
	itemsPerPage := 5
	offset := (page - 1) * itemsPerPage
	limit := itemsPerPage
	Totalitems, err := uh.CartUsecase.ExecuteTotalItemsInCart(userID)
	if Totalitems == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Empty Cart"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get the total number of products"})
		return
	}
	PaginatedCartItems, err := uh.CartUsecase.ListPaginatedCartItems(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch cartItems"})
		return
	}
	if PaginatedCartItems == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no items to show"})
		return
	}
	totalPages := (Totalitems + itemsPerPage - 1) / itemsPerPage
	response := gin.H{
		"Cart":               Cart,
		"PaginatedCartItems": PaginatedCartItems,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Cart": response})
}

// @Summary Product Delete From Cart
// @Description Removing Product From Cart
// @Tags User Cart Management
// @Accept json
// @Produce json
// @Params Cart utils.ProductDrop true ProductID
// @Success 202 "message":"product removed from your cart"
// @Failure 500 "error":"Internal Server Error"
// @Failure 400 "error":"Status Bad Request"
// @Router /deletecartproducts [delete]
func (uh *UserHandler) DeleteProductFromCart(c *gin.Context) {
	var product utils.ProductDrop
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := utils.GetUserIDFromContext(c)
	err := uh.CartUsecase.DropProduct(product.ProductID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't drop product from your cart"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Message": "product removed from your cart"})
}

//User Address Handlers

// @Summary Add Address
// @Description User Address Adding
// @Tags User Address Management
// @Accept json
// @produce json
// @Param UserAddress body entity.Address true "User address information"
// @Success 202 "Address":AddedAddress
// @Failure 500 "error":"Internal Server Error"
// @Failure 400 "error":"Status Bad Request"
// @Router /addaddress [post]
func (uh *UserHandler) AddAddress(c *gin.Context) {
	var UserAddress *entity.Address
	if err := c.ShouldBind(&UserAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := utils.GetUserIDFromContext(c)
	UserAddress.UserID = userID
	AddedAddress, err := uh.CartUsecase.AddAddress(UserAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't Add Address"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"Address": AddedAddress,
	})
}

// @Summary Edit User Address
// @Description Edit The Address By User
// @Tags User Address Management
// @Accept json
// @Produce json
// @Param UserAddress body utils.EditAddressRequest true "Edited Address"
// @Success 202 "Status Accepted":EditedAddress
// @Failure 500 "error":"Internal Server Error"
// @Failure 400 "error":"Status Bad Request"
// @Router /editaddress [patch]
func (uh *UserHandler) EditAddress(c *gin.Context) {
	var request utils.EditAddressRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := utils.GetUserIDFromContext(c)
	EditedAddress, err := uh.CartUsecase.EditAddress(request, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't edit the address"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"Edited Address": EditedAddress,
	})
}

// @Summary Update User Address
// @Description Updating Address By User
// @Tags User Address Management
// @Accept json
// @Produce json
// @Param UserAddress body utils.EditAddressRequest true "Updated Address"
// @Success 202 "Address":UpdatedAddress
// @Failure 500 "error":"Internal Server Error"
// @Failure 400 "error":"Status Bad Request"
// @Router /updateaddress [put]
func (uh *UserHandler) UpdateAddress(c *gin.Context) {
	var request utils.EditAddressRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := utils.GetUserIDFromContext(c)
	UpdatedAddress, err := uh.CartUsecase.EditAddress(request, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't edit the address"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"Address": UpdatedAddress,
	})
}

// @Summary View Addresses
// @Description View All Saved Addresses By User
// @Tags User Address Management
// @Accept json
// @Produce json
// @Success 200 "Addresses":AllAddresses
// @Failure 500 "Internal Server Error"
// @Router /viewaddresses [get]
func (uh *UserHandler) ViewAddresses(c *gin.Context) {
	userID, _ := utils.GetUserIDFromContext(c)
	AllAddresses, err := uh.CartUsecase.GetAllAddresses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't fetch address"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Addresses": AllAddresses,
	})
}

// @Summary Show Profile
// @Description Show User's Profile
// @Tags User Address Management
// @Accept json
// @Produce json
// @Success 200 "Profile":"entity.Profile","Address": entity.Address
// @Failure 500 "error":"Internal Server Error"
// @Router /showprofile [get]
func (uh *UserHandler) ShowProfile(c *gin.Context) {
	userID, _ := utils.GetUserIDFromContext(c)
	Profile, Address, err := uh.UserUsecase.FetchProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't show profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Profile": Profile,
		"Address": Address})
}

//User checkout order Handlers

// @Summary Place Order
// @Description Place Order By User
// @Tags User Checkout Management
// @Accept json
// @Produce json
// @Params Order utils.OrderRequest true "Order Request"
// @Success 202 "message":"Order Placed","Items":OrderStatus,"Address":Address
// @Failure 500 "error":"Internal Server Error"
// @Failure 400 "error":"Bad Request"
// @Router /placeorder [post]
func (uh *UserHandler) PlaceOrder(c *gin.Context) {
	var request utils.OrderRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := utils.GetUserIDFromContext(c)
	request.UserID = userID
	cart, err := uh.OrderUsecase.GetCartByUserID(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't fetch cart"})
		return
	}
	if cart.TotalPrice == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "your cart is empty"})
		return
	}
	switch request.PaymentMethod {
	case "Online":
		Order, Items, Address, err := uh.OrderUsecase.PlaceOrder(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		razorId, err := utils.Executerazorpay(cart.TotalPrice)
		fmt.Println("Id", razorId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Order Placed",
			"Order":   Order,
			"Items":   Items,
			"Address": Address,
		})
	case "COD":
		Order, Items, Address, err := uh.OrderUsecase.PlaceOrder(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't place order"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Order Placed",
			"Order":   Order,
			"Items":   Items,
			"Address": Address,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method"})
		return
	}
}

// @summary View Order
// @Description View Order By User
// @Tags User Checkout Management
// @Accept json
// @Produce json
// @Success 202 "Items":Order,"Address":Address
// @Failure 500 "error":"Internal Server Error"
// @Router /vieworder [get]
func (uh *UserHandler) ViewOrder(c *gin.Context) {
	userID, _ := utils.GetUserIDFromContext(c)
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
	TotalOrders, err := uh.OrderUsecase.ExecuteTotalOfOrders(userID)
	TotalOrderedItems, err := uh.OrderUsecase.ExecuteTotalOfOrderedItems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get the total number of products"})
		return
	}
	Address, Order, PaginatedOrderedItems, err := uh.OrderUsecase.GetOrders(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't view orders"})
		return
	}
	if len(Order) == 0 {
		c.JSON(http.StatusFound, gin.H{"message": "Currently you have no orders"})
		return
	}
	totalPages := (TotalOrderedItems + itemsPerPage - 1) / itemsPerPage
	response := gin.H{
		"Address":               Address,
		"Orders":                Order,
		"TotalOrders":           TotalOrders,
		"TotalOrderedItems":     TotalOrderedItems,
		"PaginatedOrderedItems": PaginatedOrderedItems,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Orders": response})
}

// @Summary Cancell Order
// @Description Cancell The Order By User
// @Tags User Checkout Management
// @Accept  json
// @Produce json
// @Success 202 "Order":CancelledOrder,"Address":Address
// @Failure 500 "Internal Server Error"
// @Failure 208 "message":"this item is already cancelled"
// @Failure 500 "error":"Internal Server Error"
// @Router /cancelorder [post]
func (uh *UserHandler) CancellMyOrder(c *gin.Context) {
	userID, _ := utils.GetUserIDFromContext(c)
	OrderAlreadyCancelled, err := uh.OrderUsecase.CheckOrderStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check order status"})
		return
	}
	if OrderAlreadyCancelled {
		c.JSON(http.StatusAlreadyReported, gin.H{"message": "this item is already cancelled"})
		return
	}
	Address, CancelledOrder, err := uh.OrderUsecase.CancellOrder(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't cancell order"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"Order":   CancelledOrder,
		"Address": Address,
	})
}

func (uh *UserHandler) Returnproduct(c *gin.Context) {
	var returnRequest utils.ReturnProductRequest
	if err := c.ShouldBind(&returnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uh.OrderUsecase.ExecuteReturnProduct(returnRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Request to return send successfully"})

}

//Coupon Management

func (uh *UserHandler) ShowAvailableCoupons(c *gin.Context) {
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
	TotalCoupons, err := uh.OrderUsecase.ExecuteTotalOfCoupons()
	if TotalCoupons == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no coupons available"})
		return
	}
	totalPages := (TotalCoupons + itemsPerPage - 1) / itemsPerPage
	PaginatedCoupons, err := uh.OrderUsecase.ExecutePaginatedCoupons(offset, limit)
	if PaginatedCoupons == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no coupons to show"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't show coupons"})
		return
	}
	response := gin.H{
		"PaginatedCoupons": PaginatedCoupons,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Coupons": response})
}

func (uh *UserHandler) ApplyCoupon(c *gin.Context) {
	var couponcode utils.CouponCode
	if err := c.ShouldBind(&couponcode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	code := couponcode.Code
	userID, _ := utils.GetUserIDFromContext(c)
	TotalOffer, err := uh.OrderUsecase.ExecuteApplyCoupon(code, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't apply coupon"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Offer Applied": TotalOffer})
}

//Filter Category And products

func (uh *UserHandler) FilterCategory(c *gin.Context) {
	name := c.PostForm("name")
	FilteredCategory, ProductCount, err := uh.CategoryUsecase.FetchCategoryByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "there is no category in this name"})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"Category":           FilteredCategory,
		"Number Of Products": ProductCount,
	})
}

func (uh *UserHandler) FilterProductByBrand(c *gin.Context) {
	BrandName := c.PostForm("brand_name")
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
	BrandExist, err := uh.ProductUsecase.FetchBrandName(BrandName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't verify brand name"})
		return
	}
	if BrandExist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid brand name"})
		return
	}
	TotalProducts, err := uh.ProductUsecase.ExecuteTotalOfProductsByBrand(BrandName)
	if TotalProducts == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no products in this brand"})
		return
	}
	totalPages := (TotalProducts + itemsPerPage - 1) / itemsPerPage
	PaginatedFilteredProductsByBrand, err := uh.ProductUsecase.FetchProductByBrandName(BrandName, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't show products"})
		return
	}
	if PaginatedFilteredProductsByBrand == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no products to show"})
		return
	}
	response := gin.H{
		"BrandName":               BrandName,
		"FilteredProductsByBrand": PaginatedFilteredProductsByBrand,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Products": response})
}

func (uh *UserHandler) FilterProductByItem(c *gin.Context) {
	ItemName := c.PostForm("item_name")
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
	ItemExist, err := uh.ProductUsecase.FetchItemName(ItemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't verify item name"})
		return
	}
	if ItemExist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid item name"})
		return
	}
	TotalProducts, err := uh.ProductUsecase.ExecuteTotalOfProductsByItemName(ItemName)
	if TotalProducts == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no products avalialble in this item name"})
		return
	}
	totalPages := (TotalProducts + itemsPerPage - 1) / itemsPerPage
	FilteredPaginatedProductsByItem, err := uh.ProductUsecase.FetchPaginatedProductByItemName(ItemName, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't show products"})
		return
	}
	if FilteredPaginatedProductsByItem == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no products to show"})
		return
	}
	response := gin.H{
		"ItemName":               ItemName,
		"FilteredProductsByItem": FilteredPaginatedProductsByItem,
		"Pagination": gin.H{
			"TotalPages":   totalPages,
			"CurrentPage":  page,
			"ItemsPerPage": itemsPerPage,
		},
	}
	c.JSON(http.StatusOK, gin.H{"Products": response})
}

func (uh *UserHandler) ShowWallet(c *gin.Context) {
	userID, _ := utils.GetUserIDFromContext(c)
	MyWallet, err := uh.UserUsecase.ExecuteWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't show wallet"})
		return
	}
	c.JSON(http.StatusFound, gin.H{"Wallet": MyWallet})
}

func (uh *UserHandler) VerifyMyPayment(c *gin.Context) {
	var request utils.PaymentVerificationRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	PaymentId := request.PaymentID
	OrderId := request.OrderID
	Signature := request.Signature
	err := utils.RazorPaymentVerification(Signature, OrderId, PaymentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err == nil {
		PaymentStatus := "SuccessFull"
		// OrderID, err := strconv.Atoi(OrderId)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// }
		err = uh.OrderUsecase.UpdateOrderPaymentStatus(PaymentStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't place order"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Payment Successfull"})
	}
}

func (uh *UserHandler) CreateInvoice(c *gin.Context) {
	userID, _ := utils.GetUserIDFromContext(c)
	err := uh.UserUsecase.ExecuteCreateInvoice(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't generate invoice"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "generated successfully"})
}

func (uh *UserHandler) WelcomeMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"::::": "WELCOME TO MAXPUMP"})
}
