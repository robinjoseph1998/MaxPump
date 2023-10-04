package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/utils"
)

type AdminUsecaseInterface interface {
	ExecuteAdminCreate(admin entity.Admin) (*entity.Admin, error)
	ExecuteLoginWithPassword(phone, password string) (int, error)
	ExecuteAdminLogin(phone string) error
	ExecuteOtpValidation(phone, otp string) (*entity.Admin, error)
	ExecuteAllUsersPaginated(offset int, limit int) ([]entity.User, error)
	ExecuteSearch(ID uint, FirstName string) ([]entity.User, error)
	BlockUserByID(ID uint) error
	UnBlockUserByID(ID uint) error
	CreateCoupon(coupon *entity.Coupon) (*entity.Coupon, error)
	ExecutePaginatedCoupons(offset int, limit int) ([]entity.Coupon, error)
	ExecuteTotalOfCoupons() (int, error)
	ExecuteEditCoupon(Request *utils.EditCouponRequest) (*entity.Coupon, error)
	ExecuteDeleteCoupon(code int) error
	ChangeStatus(ApplyingStatus *utils.OrderStatus) (*entity.Order, error)
	ExecuteAllOrdersPaginated(offset int, limit int) ([]entity.Order, error)
	ExecuteReturnRequestApproval(Request *utils.ReturnProductRequest) (*entity.Ordered_Item, error)
	ExecuteCreateSalesReportByDate(date string) error
	ExecuteTotalOfUsers() (int, error)
}
