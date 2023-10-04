package repository

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/utils"
	"database/sql"
	"errors"
	"fmt"

	repo "MAXPUMP1/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) repo.AdminInterface {
	return &AdminRepository{db: db}
}

func (ar *AdminRepository) Create(admin *entity.Admin) error {
	return ar.db.Create(admin).Error

}

func (ar *AdminRepository) CreateOtpKey(key, phone string) error {
	var otpkey entity.AdminOtpKey
	otpkey.Key = key
	otpkey.Phone = phone
	return ar.db.Create(&otpkey).Error
}

func (ar *AdminRepository) GetByEmail(email string) (*entity.Admin, error) {
	var admin entity.Admin
	err := ar.db.Raw("SELECT * FROM admins WHERE email = ?", email).Scan(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (ar *AdminRepository) GetByPhone(Phone string) (*entity.Admin, error) {
	var admin entity.Admin
	err := ar.db.Raw("SELECT * FROM admins WHERE phone = ?", Phone).Scan(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (ar *AdminRepository) GetAllUsersPaginated(offset int, limit int) ([]entity.User, error) {
	var users []entity.User
	err := ar.db.Raw("SELECT * FROM users OFFSET ? LIMIT ?", offset, limit).Scan(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return users, nil
}

func (ar *AdminRepository) SearchByUser(id uint, FirstName string) ([]entity.User, error) {
	var user []entity.User
	err := ar.db.Raw("SELECT * FROM users WHERE id=? AND first_name=?", id, FirstName).Scan(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (ar *AdminRepository) BlockUser(id uint) error {
	var result entity.User
	err := ar.db.Raw("SELECT * FROM users WHERE id=?", id).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	result.Blocked = true
	err = ar.db.Save(&result).Error
	if err != nil {
		return err
	}
	return nil
}

func (ar *AdminRepository) UnBlockUser(id uint) error {
	var result entity.User
	err := ar.db.Raw("SELECT * FROM users WHERE id=?", id).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	result.Blocked = false
	err = ar.db.Save(&result).Error
	if err != nil {
		return err
	}
	return nil
}

func (ar *AdminRepository) ExecuteCreateCoupon(coupon *entity.Coupon) (*entity.Coupon, error) {
	var CreatedCoupon entity.Coupon
	err := ar.db.Create(coupon).Scan(&CreatedCoupon).Error
	if err != nil {
		return nil, err
	}
	return &CreatedCoupon, nil
}

func (ar *AdminRepository) GetAllCouponsPaginated(offset int, limit int) ([]entity.Coupon, error) {
	var coupons []entity.Coupon
	err := ar.db.Raw("SELECT * FROM coupons OFFSET ? LIMIT ?", offset, limit).Scan(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

func (ar *AdminRepository) DeleteCoupon(code int) error {
	err := ar.db.Exec("DELETE FROM coupons WHERE code=?", code).Error
	if err != nil {
		return err
	}
	return nil
}

func (ar *AdminRepository) EditCoupon(FetchedCoupon *entity.Coupon) (*entity.Coupon, error) {
	var EditedCoupon entity.Coupon
	err := ar.db.Exec("UPDATE coupons SET type=?,amount=?,threshold_amount=?,expiration=?,usage_limit=? WHERE code=?",
		FetchedCoupon.Type, FetchedCoupon.Amount, FetchedCoupon.Threshold_Amount, FetchedCoupon.Expiration, FetchedCoupon.UsageLimit, FetchedCoupon.Code).Error
	if err != nil {
		return nil, err
	}
	err = ar.db.Raw("SELECT * FROM coupons WHERE code=?", FetchedCoupon.Code).Scan(&EditedCoupon).Error
	if err != nil {
		return nil, err
	}
	return &EditedCoupon, nil
}

func (ar *AdminRepository) GetCouponByCode(code int) (*entity.Coupon, error) {
	var AlreadyExistingCoupon entity.Coupon
	err := ar.db.Raw("SELECT * FROM coupons WHERE code=?", code).Scan(&AlreadyExistingCoupon).Error
	if err != nil {
		return nil, err
	}
	return &AlreadyExistingCoupon, nil
}

func (ar *AdminRepository) GetAllOrdersPaginated(offset int, limit int) ([]entity.Order, error) {
	var AllOrders []entity.Order
	err := ar.db.Raw("SELECT * FROM orders OFFSET ? LIMIT ?", offset, limit).Scan(&AllOrders).Error
	if err != nil {
		return nil, err
	}
	return AllOrders, nil
}

func (ar *AdminRepository) GetOrderByID(ID int) (*entity.Order, error) {
	var order entity.Order
	err := ar.db.Raw("SELECT * FROM orders WHERE id=?", ID).Scan(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (ar *AdminRepository) UpdateOrder(ApplyingStatus *utils.OrderStatus, Order *entity.Order) (*entity.Order, error) {
	var order entity.Order
	err := ar.db.Exec("UPDATE orders SET order_status=?,date_of_delivered=? WHERE id=?", ApplyingStatus.Status, Order.Date_Of_Delivered, ApplyingStatus.ID).Error
	if err != nil {
		return nil, err
	}
	err1 := ar.db.Raw("SELECT * FROM orders WHERE id=?", ApplyingStatus.ID).Scan(&order).Error
	if err1 != nil {
		return nil, err1
	}
	return &order, nil
}

func (ar *AdminRepository) UpdateReturnAcceptedOrder(Item *entity.Ordered_Item) (*entity.Ordered_Item, error) {
	var ApprovedItem entity.Ordered_Item
	err := ar.db.Exec("UPDATE ordered_items SET status=? WHERE product_id=?", Item.Status, Item.ProductID).Error
	if err != nil {
		return nil, err
	}
	err = ar.db.Raw("SELECT * FROM ordered_items WHERE product_id=? AND order_id=?", Item.ProductID, Item.OrderID).Scan(&ApprovedItem).Error
	if err != nil {
		return nil, err
	}
	return &ApprovedItem, nil
}

func (ar *AdminRepository) GetReturnRequestByOrderid(OrderID int, ProductID int) (*entity.ReturnedOrder, error) {
	var Request entity.ReturnedOrder
	err := ar.db.Raw("SELECT * FROM returned_orders WHERE orderid=? AND product_id=?", OrderID, ProductID).Scan(&Request).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("ReQUEST", Request)
	return &Request, nil
}

func (ar *AdminRepository) GetItemByOrderAndProductID(OrderID int, ProductID int) (*entity.Ordered_Item, error) {
	var ReturnRequestedItem entity.Ordered_Item
	err := ar.db.Raw("SELECT * FROM ordered_items WHERE order_id=? AND product_id=?", OrderID, ProductID).Scan(&ReturnRequestedItem).Error
	if err != nil {
		return nil, err
	}
	return &ReturnRequestedItem, nil
}

func (ar *AdminRepository) UpdateOrderedItem(ApplyingStatus *utils.OrderStatus) error {
	err := ar.db.Exec("UPDATE ordered_items SET status=? WHERE order_id=?", ApplyingStatus.Status, ApplyingStatus.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ar *AdminRepository) GetProductByProductID(ProductID int) (*entity.Product, error) {
	var product entity.Product
	err := ar.db.Raw("SELECT * FROM products WHERE id=?", ProductID).Scan(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (ar *AdminRepository) UpdateProductQuantity(NewQty int, ID int) error {
	err := ar.db.Exec("UPDATE products SET quantity=? WHERE id=?", NewQty, ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ar *AdminRepository) UpdateReturnRequestTable(OrderID int, ProductID int, status string) error {
	err := ar.db.Exec("UPDATE returned_orders SET status=? WHERE orderid=? AND product_id=?", status, OrderID, ProductID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ar *AdminRepository) GetSuccessfullAndDeliveredOrdersInParticularDate(Date string) ([]entity.Order, error) {
	var Orders []entity.Order
	Successfull := "SuccessFull"
	Delivered := "Delivered"
	err := ar.db.Raw("SELECT * FROM orders WHERE payment_status=? AND order_status=? AND date_of_delivered=?", Successfull, Delivered, Date).Scan(&Orders).Error
	if err != nil {
		return nil, err
	}
	return Orders, nil
}

func (ar *AdminRepository) GetTotalAmountOfSuccessfullAndDeliveredOrders(date string) (float64, error) {
	var Total sql.NullFloat64 // Using sql.NullFloat64 to handle NULL values
	Successfull := "SuccessFull"
	Delivered := "Delivered"
	err := ar.db.Raw("SELECT SUM(total_price) FROM orders WHERE payment_status=? AND order_status=? AND date_of_delivered=?", Successfull, Delivered, date).Scan(&Total).Error
	if err != nil {
		return 0, err
	}
	if !Total.Valid {
		return 0, nil
	}
	return Total.Float64, nil
}

func (ar *AdminRepository) GetOrderedItemsByOrderID(OrdersID int) ([]entity.Ordered_Item, error) {
	var Items []entity.Ordered_Item
	err := ar.db.Raw("SELECT * FROM ordered_items WHERE order_id=?", OrdersID).Scan(&Items).Error
	if err != nil {
		return nil, err
	}
	return Items, nil
}

func (ar *AdminRepository) GetTotalOfUsers() (int, error) {
	var count int
	err := ar.db.Raw("SELECT COUNT(*) FROM users").Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ar *AdminRepository) GetTotalOfCoupons() (int, error) {
	var count int
	err := ar.db.Raw("SELECT COUNT(*)FROM coupons").Scan(&count).Error
	if err != nil {
		return 0, nil
	}
	return count, nil
}
