package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/utils"
	"time"
)

type OrderUsecaseInterface interface {
	PlaceOrder(request utils.OrderRequest) (*entity.Order, []entity.Ordered_Item, *entity.Address, error)
	GetOrders(userID int, offset int, limit int) ([]entity.Address, []entity.Order, []entity.Ordered_Item, error)
	CancellOrder(userID int) (*entity.Address, *entity.Order, error)
	CheckOrderStatus(userID int) (bool, error)
	ExecutePaginatedCoupons(offset int, limit int) ([]entity.Coupon, error)
	ExecuteApplyCoupon(code int, userID int) (*entity.Cart, error)
	GetCartByUserID(userID int) (*entity.Cart, error)
	ExecuteReturnProduct(returnRequest utils.ReturnProductRequest) error
	ExecuteShowReturnRequestsPaginated(offset int, limit int) (*entity.ReturnedOrder, error)
	UpdateOrderPaymentStatus(PaymentStatus string) error
	ExecuteSuccessfullOrdersPaginated(offset int, limit int) ([]entity.Order, error)
	ExecuteSalesInParticularDatePaginated(ParsedTime time.Time, offset int, limit int) ([]entity.Order, error)
	ExecuteTotalOfOrders(userID int) (int, error)
	ExecuteTotalOfOrderedItems(userID int) (int, error)
	ExecuteTotalOfCoupons() (int, error)
	ExecuteTotalOfSales() (int, error)
	ExecuteTotalOfSalesByDate(Date time.Time) (int, error)
	ExecuteTotalOfAllOrders() (int, error)
	ExecuteTotalOfReturnRequests() (int, error)
}
