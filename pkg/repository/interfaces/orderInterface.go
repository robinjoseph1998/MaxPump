package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"time"
)

type OrderInterface interface {
	GetCartByUserID(ID int) (*entity.Cart, error)
	GetAddressByUserID(userID int) (*entity.Address, error)
	GetAddressesByUserID(userID int) ([]entity.Address, error)
	OrderCreater(Order *entity.Order) (*entity.Order, error)
	DeleteCartByUserID(userID int) error
	GetOrdersByUserID(userID int) ([]entity.Order, error)
	ChangeOrderStatus(userID int, status string) (*entity.Order, error)
	CheckStatus(userID int) (bool, error)
	GetPaginatedCoupons(offset int, limit int) ([]entity.Coupon, error)
	GetCouponByID(code int) (*entity.Coupon, error)
	ApplyCoupon(discount float64, userID int) (*entity.Cart, error)
	UpdateCoupon(UsedCount int, code int) error
	GetCartItemsByUserID(userID int) ([]entity.CartItem, error)
	DeleteCartItemsByCartID(cartID int) error
	DecreaseQuantityOfProducts(UpdatedQuantity int, ProductID int) error
	GetCartItemsByProductID(ProductID int) (*entity.CartItem, error)
	GetProductsByID(ID int) (*entity.Product, error)
	GetCartItemsByIDofUser(userID int) (*entity.CartItem, error)
	GetOrderByOrderID(ID int) (*entity.Order, error)
	SendToRequestTable(request *entity.ReturnedOrder) (*entity.ReturnedOrder, error)
	GetReturnRequestTablePaginated(offset int, limit int) (*entity.ReturnedOrder, error)
	UpdateOrderedItem(Item *entity.Ordered_Item) error
	FetchOrderedItems(OrderID int) ([]entity.Ordered_Item, error)
	GetItemByOrderAndProductID(OrderID int, ProductID int) (*entity.Ordered_Item, error)
	FetchOrderedItemsByUserID(userID int, offset int, limit int) ([]entity.Ordered_Item, error)
	UpdateUserWallet(UserWallet *entity.Wallet) error
	EditOrderPaymentStatus(PaymentStatus string) error
	GetSuccessfullOrdersPaginated(offset int, limit int) ([]entity.Order, error)
	GetSuccessfullOrderByDate(ParsedTime time.Time, offset int, limit int) ([]entity.Order, error)
	GetTotalOrdersOfUser(userID int) (int, error)
	GetTotalOrderedItemsOfUser(userID int) (int, error)
	GetTotalOfCoupons() (int, error)
	GetTotalOfSales() (int, error)
	GetTotalOfSalesByDate(Date time.Time) (int, error)
	GetTotalOfAllOrders() (int, error)
	GetTotalOfReturnRequests() (int, error)
}
