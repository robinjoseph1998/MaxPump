package repository

import (
	"MAXPUMP1/pkg/domain/entity"
	repo "MAXPUMP1/pkg/repository/interfaces"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repo.OrderInterface {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) GetCartByUserID(ID int) (*entity.Cart, error) {
	var cart entity.Cart
	err := or.db.Raw("SELECT * FROM carts WHERE user_id=?", ID).Scan(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (or *OrderRepository) GetAddressByUserID(userID int) (*entity.Address, error) {
	var address entity.Address
	err := or.db.Raw("SELECT * FROM addresses WHERE user_id=?", userID).Scan(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &address, nil

}

func (or *OrderRepository) OrderCreater(Order *entity.Order) (*entity.Order, error) {
	fmt.Println("CREATER TIME", Order.Date_Of_Delivered)
	var CreatedOrder entity.Order
	err := or.db.Create(Order).Scan(&CreatedOrder).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("ORDER TIME", CreatedOrder.Date_Of_Delivered)
	return &CreatedOrder, nil

}

func (or *OrderRepository) DeleteCartByUserID(userID int) error {
	err := or.db.Exec("DELETE FROM carts WHERE user_id=?", userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) GetCartItemsByUserID(userID int) ([]entity.CartItem, error) {
	var Items []entity.CartItem
	err := or.db.Raw("SELECT * FROM cart_items WHERE user_id=?", userID).Scan(&Items).Error
	if err != nil {
		return nil, err
	}
	return Items, nil
}

func (or *OrderRepository) GetOrdersByUserID(userID int) ([]entity.Order, error) {
	var Orders []entity.Order
	err := or.db.Raw("SELECT * FROM orders WHERE user_id=?", userID).Scan(&Orders).Error
	if err != nil {
		return nil, err
	}
	return Orders, nil
}

func (or *OrderRepository) DeleteCartItemsByCartID(cartID int) error {
	err := or.db.Exec("DELETE FROM cart_items WHERE cart_id=?", cartID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) GetAddressesByUserID(userID int) ([]entity.Address, error) {
	var addresses []entity.Address
	err := or.db.Raw("SELECT * FROM addresses WHERE user_id=?", userID).Scan(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (or *OrderRepository) ChangeOrderStatus(userID int, status string) (*entity.Order, error) {
	var CancelledOrder entity.Order
	err := or.db.Raw("UPDATE orders SET order_status=? WHERE user_id=?", status, userID).Scan(&CancelledOrder).Error
	if err != nil {
		return nil, err
	}
	return &CancelledOrder, nil
}

func (or *OrderRepository) CheckStatus(userID int) (bool, error) {
	var orders entity.Order
	err := or.db.Raw("SELECT * FROM orders WHERE user_id=?", userID).Scan(&orders).Error
	if err != nil {
		return false, err
	}
	if orders.OrderStatus == "cancelled" {
		return true, nil
	}
	return false, nil
}

func (or *OrderRepository) GetPaginatedCoupons(offset int, limit int) ([]entity.Coupon, error) {
	var Coupons []entity.Coupon
	err := or.db.Raw("SELECT * FROM coupons OFFSET ? LIMIT ?", offset, limit).Scan(&Coupons).Error
	if err != nil {
		return nil, err
	}
	return Coupons, nil
}

func (or *OrderRepository) GetCouponByID(code int) (*entity.Coupon, error) {
	var coupon entity.Coupon
	err := or.db.Raw("SELECT * FROM coupons WHERE code=?", code).Scan(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (or *OrderRepository) ApplyCoupon(discount float64, userID int) (*entity.Cart, error) {
	var cart entity.Cart
	err := or.db.Exec("UPDATE carts SET total_price=? WHERE user_id=?", discount, userID).Error
	if err != nil {
		return nil, err
	}
	err1 := or.db.Raw("SELECT * FROM carts WHERE user_id=?", userID).Scan(&cart).Error
	if err1 != nil {
		return nil, err1
	}
	return &cart, nil
}

func (or *OrderRepository) UpdateCoupon(UsedCount int, code int) error {
	err := or.db.Exec("UPDATE coupons SET used_count=? WHERE code=?", UsedCount, code).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) DecreaseQuantityOfProducts(quantity int, ProductID int) error {
	err := or.db.Exec("UPDATE products SET quantity=? WHERE id=?", quantity, ProductID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) GetCartItemsByProductID(ProductID int) (*entity.CartItem, error) {
	var items entity.CartItem
	err := or.db.Raw("SELECT * FROM cart_items WHERE product_item_id=?", ProductID).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func (or *OrderRepository) GetProductsByID(ID int) (*entity.Product, error) {
	var product entity.Product
	err := or.db.Raw("SELECT * FROM products WHERE id=?", ID).Scan(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (or *OrderRepository) GetCartItemsByIDofUser(userID int) (*entity.CartItem, error) {
	var Items entity.CartItem
	err := or.db.Raw("SELECT * FROM cart_items WHERE user_id=?", userID).Scan(&Items).Error
	if err != nil {
		return nil, err
	}
	return &Items, nil
}

func (or *OrderRepository) GetOrderByOrderID(ID int) (*entity.Order, error) {
	var Order entity.Order
	err := or.db.Raw("SELECT * FROM orders WHERE id=?", ID).Scan(&Order).Error
	if err != nil {
		return nil, err
	}
	return &Order, nil
}
func (or *OrderRepository) SendToRequestTable(request *entity.ReturnedOrder) (*entity.ReturnedOrder, error) {
	var Product entity.ReturnedOrder
	err := or.db.Create(request).Scan(&Product).Error
	if err != nil {
		return nil, err
	}
	return &Product, nil
}

func (or *OrderRepository) GetReturnRequestTablePaginated(offset int, limit int) (*entity.ReturnedOrder, error) {
	var Requests entity.ReturnedOrder
	err := or.db.Raw("SELECT * FROM returned_orders OFFSET ? LIMIT ?", offset, limit).Scan(&Requests).Error
	if err != nil {
		return nil, err
	}
	return &Requests, nil
}

func (or *OrderRepository) UpdateOrderedItem(Item *entity.Ordered_Item) error {
	err := or.db.Create(Item).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) FetchOrderedItems(OrderID int) ([]entity.Ordered_Item, error) {
	var Items []entity.Ordered_Item
	err := or.db.Raw("SELECT * FROM ordered_items WHERE order_id=?", OrderID).Scan(&Items).Error
	if err != nil {
		return nil, err
	}
	return Items, nil
}

func (or *OrderRepository) GetItemByOrderAndProductID(OrderID int, ProductID int) (*entity.Ordered_Item, error) {
	var item entity.Ordered_Item
	err := or.db.Raw("SELECT * FROM ordered_items WHERE order_id=? AND product_id=?", OrderID, ProductID).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (or *OrderRepository) FetchOrderedItemsByUserID(userID int, offset int, limit int) ([]entity.Ordered_Item, error) {
	var items []entity.Ordered_Item
	err := or.db.Raw("SELECT * FROM ordered_items WHERE user_id=? OFFSET ? LIMIT ?", userID, offset, limit).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (or *OrderRepository) UpdateUserWallet(UserWallet *entity.Wallet) error {
	err := or.db.Create(UserWallet).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) EditOrderPaymentStatus(PaymentStatus string) error {
	err := or.db.Exec("UPDATE orders SET payment_status=?", PaymentStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) GetSuccessfullOrdersPaginated(offset int, limit int) ([]entity.Order, error) {
	var SuccessFullOrders []entity.Order
	Successfull := "SuccessFull"
	Delivered := "Delivered"
	err := or.db.Raw("SELECT * FROM orders WHERE payment_status=? AND order_status=? OFFSET ? LIMIT ?", Successfull, Delivered, offset, limit).Scan(&SuccessFullOrders).Error
	if err != nil {
		return nil, err
	}
	return SuccessFullOrders, nil
}

func (or *OrderRepository) GetSuccessfullOrderByDate(ParsedTime time.Time, offset int, limit int) ([]entity.Order, error) {
	var SuccessfullSalesInThisDate []entity.Order
	Successfull := "SuccessFull"
	Delivered := "Delivered"
	err := or.db.Raw("SELECT * FROM orders WHERE payment_status=? AND order_status=? AND DATE(date_of_delivered)=? OFFSET ? LIMIT ?", Successfull, Delivered, ParsedTime, offset, limit).Scan(&SuccessfullSalesInThisDate).Error
	if err != nil {
		return nil, err
	}
	return SuccessfullSalesInThisDate, nil
}

func (or *OrderRepository) GetTotalOrdersOfUser(userID int) (int, error) {
	var count int
	err := or.db.Raw("SELECT COUNT(*) FROM orders WHERE user_id=?", userID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (or *OrderRepository) GetTotalOrderedItemsOfUser(userID int) (int, error) {
	var count int
	err := or.db.Raw("SELECT COUNT(*) FROM ordered_items WHERE user_id=?", userID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (or *OrderRepository) GetTotalOfCoupons() (int, error) {
	var count int
	err := or.db.Raw("SELECT COUNT(*)FROM coupons").Scan(&count).Error
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func (or *OrderRepository) GetTotalOfSales() (int, error) {
	var count int
	Successfull := "SuccessFull"
	Delivered := "Delivered"
	err := or.db.Raw("SELECT COUNT(*) FROM orders WHERE payment_status=? AND order_status=?", Successfull, Delivered).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (or *OrderRepository) GetTotalOfSalesByDate(Date time.Time) (int, error) {
	var count int
	Successfull := "SuccessFull"
	Delivered := "Delivered"
	err := or.db.Raw("SELECT COUNT(*) FROM orders WHERE payment_status=? AND order_status=? AND DATE(date_of_delivered)=?", Successfull, Delivered, Date).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (or *OrderRepository) GetTotalOfAllOrders() (int, error) {
	var count int
	err := or.db.Raw("SELECT COUNT(*) FROM orders").Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (or *OrderRepository) GetTotalOfReturnRequests() (int, error) {
	var count int
	err := or.db.Raw("SELECT COUNT(*) FROM returned_orders").Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
