package usecase

import (
	"MAXPUMP1/pkg/domain/entity"
	repo "MAXPUMP1/pkg/repository/interfaces"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"MAXPUMP1/pkg/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	adminRepo repo.AdminInterface
}

func NewAdmin(adminRepo repo.AdminInterface) use.AdminUsecaseInterface {
	return &AdminUsecase{adminRepo: adminRepo}

}

func (au *AdminUsecase) ExecuteAdminCreate(admin entity.Admin) (*entity.Admin, error) {
	email, err := au.adminRepo.GetByEmail(admin.Email)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if email.Email != "" {
		return nil, errors.New("admin with this email already exists")
	}
	phone, err := au.adminRepo.GetByPhone(admin.Phone)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if phone.Phone != "" {
		return nil, errors.New("user with this phone no already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newAdmin := &entity.Admin{
		AdminName: admin.AdminName,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Password:  string(hashedPassword),
	}
	err1 := au.adminRepo.Create(newAdmin)
	if err1 != nil {
		return nil, err1
	}
	return newAdmin, nil
}

func (au *AdminUsecase) ExecuteLoginWithPassword(phone, password string) (int, error) {
	admin, err := au.adminRepo.GetByPhone(phone)
	if err != nil {
		return 0, err
	}
	if admin == nil {
		return 0, errors.New("admin with this phone number not exist")
	}
	return int(admin.ID), nil
}

func (au *AdminUsecase) ExecuteAdminLogin(phone string) error {
	result, err := au.adminRepo.GetByPhone(phone)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("admin with this phone not exist")
	}
	key, err1 := utils.SendOtp(phone)
	if err1 != nil {
		return err1
	} else {
		err := au.adminRepo.CreateOtpKey(key, phone)
		if err != nil {
			return err
		}
	}
	return nil
}

func (au *AdminUsecase) ExecuteOtpValidation(phone, otp string) (*entity.Admin, error) {
	result, err := au.adminRepo.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	err1 := utils.CheckOtp(phone, otp)
	if err1 != nil {
		return nil, err1
	}
	return result, nil
}

func (au *AdminUsecase) ExecuteAllUsersPaginated(offset int, limit int) ([]entity.User, error) {
	users, err := au.adminRepo.GetAllUsersPaginated(offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (au *AdminUsecase) ExecuteSearch(ID uint, FirstName string) ([]entity.User, error) {
	user, err := au.adminRepo.SearchByUser(ID, FirstName)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (au *AdminUsecase) BlockUserByID(ID uint) error {
	err := au.adminRepo.BlockUser(ID)
	if err != nil {
		return err
	}
	return nil

}

func (au *AdminUsecase) UnBlockUserByID(ID uint) error {
	err := au.adminRepo.UnBlockUser(ID)
	if err != nil {
		return err
	}
	return nil
}

func (au *AdminUsecase) CreateCoupon(coupon *entity.Coupon) (*entity.Coupon, error) {
	ExistingCoupon, err := au.adminRepo.GetCouponByCode(coupon.Code)
	if err != nil {
		return nil, errors.New("can't get coupon ")
	}
	if ExistingCoupon.Code == coupon.Code {
		return nil, errors.New("the coupon with same code already exists")
	}
	CreatedCoupon, err := au.adminRepo.ExecuteCreateCoupon(coupon)
	if err != nil {
		return nil, err
	}
	return CreatedCoupon, nil
}

func (au *AdminUsecase) ExecutePaginatedCoupons(offset int, limit int) ([]entity.Coupon, error) {
	AllCoupons, err := au.adminRepo.GetAllCouponsPaginated(offset, limit)
	if err != nil {
		return nil, err
	}
	return AllCoupons, nil
}

func (au *AdminUsecase) ExecuteEditCoupon(Request *utils.EditCouponRequest) (*entity.Coupon, error) {
	var EditedCoupon *entity.Coupon
	FetchedCoupon, err := au.adminRepo.GetCouponByCode(Request.Code)
	if err != nil {
		return nil, errors.New("can't fetch the coupon in database")
	}
	if FetchedCoupon.Type == "" {
		return nil, errors.New("invalid coupon code")
	}
	if Request.Amount != nil {
		FetchedCoupon.Amount = *Request.Amount
	}
	if Request.Threshold_Amount != nil {
		FetchedCoupon.Threshold_Amount = *Request.Threshold_Amount
	}
	if Request.Expiration != nil {
		FetchedCoupon.Expiration = *Request.Expiration
	}
	if Request.UsageLimit != nil {
		FetchedCoupon.UsageLimit = *Request.UsageLimit
	}
	EditedCoupon, err = au.adminRepo.EditCoupon(FetchedCoupon)
	if err != nil {
		return nil, errors.New("can't edit coupon")
	}
	return EditedCoupon, nil
}

func (au *AdminUsecase) ExecuteDeleteCoupon(code int) error {
	if err := au.adminRepo.DeleteCoupon(code); err != nil {
		return err
	}
	return nil
}

func (au *AdminUsecase) ChangeStatus(ApplyingStatus *utils.OrderStatus) (*entity.Order, error) {
	Order, err := au.adminRepo.GetOrderByID(ApplyingStatus.ID)
	if err != nil {
		return nil, err
	}
	if Order.OrderStatus == ApplyingStatus.Status {
		return nil, errors.New("this status is already updated in this order")
	}
	if ApplyingStatus.Status == "Delivered" {
		layout := "2006-01-02"
		date := time.Now()
		formattedDate := date.Format(layout)
		Order.Date_Of_Delivered = formattedDate
		Order.PaymentStatus = "Successfull"
	} else {
		Order.Date_Of_Delivered = ""
		Order.PaymentStatus = "Pending"
	}
	StatusUpdatedOrder, err := au.adminRepo.UpdateOrder(ApplyingStatus, Order)
	if err != nil {
		return nil, err
	}
	err = au.adminRepo.UpdateOrderedItem(ApplyingStatus)
	if err != nil {
		return nil, errors.New("can't update ordered item")
	}
	return StatusUpdatedOrder, nil
}

func (au *AdminUsecase) ExecuteAllOrdersPaginated(offset int, limit int) ([]entity.Order, error) {
	AllOrders, err := au.adminRepo.GetAllOrdersPaginated(offset, limit)
	if err != nil {
		return nil, errors.New("can't fetch orders")
	}
	return AllOrders, nil
}

func (au *AdminUsecase) ExecuteReturnRequestApproval(Request *utils.ReturnProductRequest) (*entity.Ordered_Item, error) {
	ReturnRequests, err := au.adminRepo.GetReturnRequestByOrderid(Request.OrderID, Request.ProductID)
	if err != nil {
		return nil, errors.New("can't fetch return requests for verify")
	}
	if ReturnRequests.Status == "Returned" {
		return nil, errors.New("this request is already approved")
	}
	Item, err := au.adminRepo.GetItemByOrderAndProductID(Request.OrderID, Request.ProductID)
	if err != nil {
		return nil, errors.New("can't fetch order")
	}
	Item.Status = "Returned"
	ApprovedOrderedItem, err := au.adminRepo.UpdateReturnAcceptedOrder(Item)
	if err != nil {
		return nil, errors.New("can't update order")
	}
	if ApprovedOrderedItem.Status == "Returned" {
		product, err := au.adminRepo.GetProductByProductID(ApprovedOrderedItem.ProductID)
		if err != nil {
			return nil, errors.New("can't fetch product for quantity updation")
		}
		NewQty := product.Quantity + ApprovedOrderedItem.Quantity
		err = au.adminRepo.UpdateProductQuantity(NewQty, ApprovedOrderedItem.ProductID)
		if err != nil {
			return nil, errors.New("can't update quantity in product base")
		}
		err = au.adminRepo.UpdateReturnRequestTable(Request.OrderID, Request.ProductID, Item.Status)
		if err != nil {
			return nil, errors.New("can't update return request table")
		}
	}
	return ApprovedOrderedItem, nil
}

func (au *AdminUsecase) ExecuteCreateSalesReportByDate(date string) error {
	Orders, err := au.adminRepo.GetSuccessfullAndDeliveredOrdersInParticularDate(date)
	if err != nil {
		return errors.New("can't fetch orders")
	}
	SalesCount := len(Orders)
	TotalAmount, err := au.adminRepo.GetTotalAmountOfSuccessfullAndDeliveredOrders(date)
	if err != nil {
		return errors.New("can't fetch total amount of successfull orders")
	}
	SalesByProducts := make(map[string]int)
	PricePerQuantity := make(map[string]float64)
	for _, order := range Orders {
		OrderedItems, err := au.adminRepo.GetOrderedItemsByOrderID(order.ID)
		if err != nil {
			return errors.New("can't get products by sales")
		}
		for _, item := range OrderedItems {
			productName := item.Item
			quantity := item.Quantity
			PricePerUnit := item.Price * float64(quantity)
			if _, exists := SalesByProducts[productName]; exists {
				SalesByProducts[productName] += quantity
				PricePerQuantity[productName] += PricePerUnit
			} else {
				SalesByProducts[productName] = quantity
				PricePerQuantity[productName] = PricePerUnit
			}
		}
	}
	err = utils.SalesReportGenerator(date, SalesCount, TotalAmount, SalesByProducts, PricePerQuantity)
	if err != nil {
		return errors.New("can't create sales report")
	}
	return nil
}

func (au *AdminUsecase) ExecuteTotalOfUsers() (int, error) {
	totalUsers, err := au.adminRepo.GetTotalOfUsers()
	if err != nil {
		return 0, errors.New("can't fetch the total count of the users")
	}
	return totalUsers, nil
}

func (au *AdminUsecase) ExecuteTotalOfCoupons() (int, error) {
	TotalOfAvaialableCoupons, err := au.adminRepo.GetTotalOfCoupons()
	if err != nil {
		return 0, errors.New("can't fetch total of the available coupons")
	}
	return TotalOfAvaialableCoupons, nil
}
