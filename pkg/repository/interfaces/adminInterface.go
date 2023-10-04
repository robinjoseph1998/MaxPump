package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/utils"
)

type AdminInterface interface {

	//Admin creation and Login interfaces
	Create(admin *entity.Admin) error
	CreateOtpKey(key, phone string) error
	GetByPhone(phone string) (*entity.Admin, error)
	GetByEmail(email string) (*entity.Admin, error)

	//Admin user Management Interfaces
	GetAllUsersPaginated(offset int, limit int) ([]entity.User, error)
	SearchByUser(id uint, FirstName string) ([]entity.User, error)
	BlockUser(id uint) error
	UnBlockUser(id uint) error
	GetTotalOfUsers() (int, error)

	//Admin Coupon Management Interfaces
	ExecuteCreateCoupon(coupon *entity.Coupon) (*entity.Coupon, error)
	GetAllCouponsPaginated(offset int, limit int) ([]entity.Coupon, error)
	GetTotalOfCoupons() (int, error)
	DeleteCoupon(code int) error
	GetCouponByCode(code int) (*entity.Coupon, error)
	EditCoupon(FetchedCoupon *entity.Coupon) (*entity.Coupon, error)

	//Admin Order Management
	GetAllOrdersPaginated(offset int, limit int) ([]entity.Order, error)
	GetOrderByID(ID int) (*entity.Order, error)
	UpdateOrder(ApplyingStatus *utils.OrderStatus, Order *entity.Order) (*entity.Order, error)
	UpdateReturnAcceptedOrder(Item *entity.Ordered_Item) (*entity.Ordered_Item, error)
	GetReturnRequestByOrderid(OrderID int, ProductID int) (*entity.ReturnedOrder, error)
	GetItemByOrderAndProductID(OrderID int, ProductID int) (*entity.Ordered_Item, error)
	UpdateOrderedItem(ApplyingStatus *utils.OrderStatus) error
	GetProductByProductID(ProductID int) (*entity.Product, error)
	UpdateProductQuantity(NewQty int, ID int) error
	UpdateReturnRequestTable(OrderID int, ProductID int, status string) error
	GetSuccessfullAndDeliveredOrdersInParticularDate(Date string) ([]entity.Order, error)
	GetTotalAmountOfSuccessfullAndDeliveredOrders(date string) (float64, error)
	GetOrderedItemsByOrderID(OrdersID int) ([]entity.Ordered_Item, error)
}
