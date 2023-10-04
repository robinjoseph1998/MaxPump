package usecase

import (
	"MAXPUMP1/pkg/domain/entity"
	repo "MAXPUMP1/pkg/repository/interfaces"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"MAXPUMP1/pkg/utils"
	"errors"
	"time"
)

type OrderUsecase struct {
	orderRepo repo.OrderInterface
}

func NewOrder(orderRepo repo.OrderInterface) use.OrderUsecaseInterface {
	return &OrderUsecase{orderRepo: orderRepo}

}

func (ou *OrderUsecase) PlaceOrder(request utils.OrderRequest) (*entity.Order, []entity.Ordered_Item, *entity.Address, error) {
	var Order *entity.Order
	cart, err := ou.orderRepo.GetCartByUserID(request.UserID)
	if err != nil {
		return nil, nil, nil, errors.New("can't fetch cart")
	}
	if cart.UserID != 0 {
		err := ou.orderRepo.DeleteCartByUserID(cart.UserID)
		if err != nil {
			return nil, nil, nil, errors.New("can't removed ordered cart")
		}
	}
	address, err := ou.orderRepo.GetAddressByUserID(request.AddressId)
	if err != nil {
		return nil, nil, nil, errors.New("can't fetch address")
	}
	if address.UserID != request.UserID {
		return nil, nil, nil, errors.New("user Id not matching")
	}
	Order = &entity.Order{
		UserID:            address.UserID,
		AddressId:         int(address.ID),
		PaymentMethod:     request.PaymentMethod,
		TotalPrice:        cart.TotalPrice,
		OrderStatus:       "Processing",
		PaymentStatus:     "Pending",
		Date_Of_Delivered: "Not Delivered Yet",
	}
	CreatedOrder, err := ou.orderRepo.OrderCreater(Order)
	if err != nil {
		return nil, nil, nil, err
	}
	if CreatedOrder.OrderStatus == "Processing" {
		cart_items, err := ou.orderRepo.GetCartItemsByUserID(request.UserID)
		if err != nil {
			return nil, nil, nil, errors.New("can't fetch cart items")
		}
		for _, cartItem := range cart_items {
			product, err := ou.orderRepo.GetProductsByID(cartItem.ProductItemID)
			if err != nil {
				return nil, nil, nil, err
			}
			UpdatedQuantity := product.Quantity - cartItem.Qty
			if UpdatedQuantity < 0 {
				return nil, nil, nil, errors.New("Insufficient quantity of product: " + product.Brand_Name + product.Description)
			}
			err1 := ou.orderRepo.DecreaseQuantityOfProducts(UpdatedQuantity, int(product.ID))
			if err1 != nil {
				return nil, nil, nil, errors.New("can't decrease the quantity of products from products database")
			}
			Ordered_Item := &entity.Ordered_Item{
				ProductID: int(product.ID),
				Item:      product.Item,
				Quantity:  cartItem.Qty,
				Price:     product.Price,
				UserID:    CreatedOrder.UserID,
				OrderID:   CreatedOrder.ID,
				Status:    CreatedOrder.OrderStatus,
			}
			err = ou.orderRepo.UpdateOrderedItem(Ordered_Item)
			if err != nil {
				return nil, nil, nil, errors.New("can't update ordered item")
			}
		}
		CartItems, err := ou.orderRepo.GetCartItemsByIDofUser(request.UserID)
		if err != nil {
			return nil, nil, nil, errors.New("can't fetch cartitems for deletion")
		}
		if CartItems.CartId != 0 {
			err := ou.orderRepo.DeleteCartItemsByCartID(CartItems.CartId)
			if err != nil {
				return nil, nil, nil, errors.New("can't removed ordered cart_items")
			}
		}
	}
	Ordered_Items, err := ou.orderRepo.FetchOrderedItems(CreatedOrder.ID)
	if err != nil {
		return nil, nil, nil, errors.New("can't fetch ordered items")
	}
	return CreatedOrder, Ordered_Items, address, nil
}

func (ou *OrderUsecase) UpdateOrderPaymentStatus(PaymentStatus string) error {
	err := ou.orderRepo.EditOrderPaymentStatus(PaymentStatus)
	if err != nil {
		return err
	}
	return nil
}

func (ou *OrderUsecase) GetOrders(userID int, offset int, limit int) ([]entity.Address, []entity.Order, []entity.Ordered_Item, error) {
	address, err := ou.orderRepo.GetAddressesByUserID(userID)
	if err != nil {
		return nil, nil, nil, errors.New("can't fetch address")
	}
	order, err := ou.orderRepo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, nil, nil, errors.New("can't fetch address")
	}
	items, err := ou.orderRepo.FetchOrderedItemsByUserID(userID, offset, limit)
	if err != nil {
		return nil, nil, nil, errors.New("can't fetch ordered items")
	}
	return address, order, items, nil
}

func (ou *OrderUsecase) CancellOrder(userID int) (*entity.Address, *entity.Order, error) {
	status := "cancelled"
	address, err := ou.orderRepo.GetAddressByUserID(userID)
	if err != nil {
		return nil, nil, errors.New("cant find fetch address")
	}
	order, err := ou.orderRepo.ChangeOrderStatus(userID, status)
	if err != nil {
		return nil, nil, errors.New("can't change order status")
	}
	return address, order, nil
}

func (ou *OrderUsecase) CheckOrderStatus(userID int) (bool, error) {
	ItemAlreadyCancelled, err := ou.orderRepo.CheckStatus(userID)
	if err != nil {
		return false, err
	}
	if ItemAlreadyCancelled {
		return true, nil
	}
	return false, nil
}

func (ou *OrderUsecase) ExecutePaginatedCoupons(offset int, limit int) ([]entity.Coupon, error) {
	AllCoupons, err := ou.orderRepo.GetPaginatedCoupons(offset, limit)
	if err != nil {
		return nil, err
	}
	return AllCoupons, nil
}

func (ou *OrderUsecase) ExecuteApplyCoupon(code int, userID int) (*entity.Cart, error) {
	coupon, err := ou.orderRepo.GetCouponByID(code)
	if err != nil {
		return nil, errors.New("can't fetch coupon")
	}
	currentTime := time.Now()
	if currentTime.Before(coupon.Expiration) {
		return nil, errors.New("this coupon is expired")
	}
	cart, err := ou.orderRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cant fetch cart")
	}
	if cart.TotalPrice < coupon.Threshold_Amount {
		return nil, errors.New("the cart total is lower than the threshold price for apply this coupon")
	}
	if coupon.UsedCount == coupon.UsageLimit {
		return nil, errors.New("coupon limit exceed")
	}
	discount := cart.TotalPrice - coupon.Amount
	CouponAppliedCart, err := ou.orderRepo.ApplyCoupon(discount, userID)
	if err != nil {
		return nil, errors.New("can't apply coupon")
	}
	var UsedCount int
	UsedCount += 1
	err1 := ou.orderRepo.UpdateCoupon(UsedCount, code)
	if err1 != nil {
		return nil, errors.New("can't update coupon")
	}
	return CouponAppliedCart, nil
}

func (ou *OrderUsecase) GetCartByUserID(userID int) (*entity.Cart, error) {
	MyCart, err := ou.orderRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("can't fetch cart")
	}
	return MyCart, nil
}

func (ou *OrderUsecase) ExecuteReturnProduct(returnRequest utils.ReturnProductRequest) error {
	var ReturnProductRequest entity.ReturnedOrder
	Order, err := ou.orderRepo.GetOrderByOrderID(returnRequest.OrderID)
	if err != nil {
		return errors.New("can't fetch order")
	}
	Ordered_Item, err := ou.orderRepo.GetItemByOrderAndProductID(returnRequest.OrderID, returnRequest.ProductID)
	if err != nil {
		return errors.New("can't fetch Ordered Item")
	}
	currentDate := time.Now()
	deliveredDate, err := time.Parse("2006-1-2", Order.Date_Of_Delivered)
	if err != nil {
		return errors.New("error in converting the type of date")
	}
	timeDifference := currentDate.Sub(deliveredDate)
	const (
		StatusOfProduct   = "Delivered"
		MaxReturnDuration = 10 * 24 * time.Hour
	)
	switch Order.OrderStatus {
	case StatusOfProduct:
		if timeDifference < MaxReturnDuration {
			ReturnProductRequest.Price = float64(Ordered_Item.Price)
			ReturnProductRequest.UserID = Order.UserID
			ReturnProductRequest.Orderid = Order.ID
			ReturnProductRequest.ProductID = Ordered_Item.ProductID
			ReturnProductRequest.Item = Ordered_Item.Item
			ReturnProductRequest.Date_Of_Delivered = Order.Date_Of_Delivered
			ReturnProductRequest.Status = Order.OrderStatus
			ReturnedProduct, err := ou.orderRepo.SendToRequestTable(&ReturnProductRequest)
			if err != nil {
				return errors.New("can't send request to return")
			}
			UserWallet := entity.Wallet{
				UserID:   ReturnedProduct.UserID,
				Credited: ReturnedProduct.Price,
				Balance:  ReturnedProduct.Price,
			}
			err = ou.orderRepo.UpdateUserWallet(&UserWallet)
			if err != nil {
				return errors.New("can't update wallet")
			}
		} else {
			return errors.New("return period exceeded")
		}
	default:
		return errors.New("invalid option")
	}
	return nil
}

func (ou *OrderUsecase) ExecuteShowReturnRequestsPaginated(offset int, limit int) (*entity.ReturnedOrder, error) {
	ReturnedRequestTable, err := ou.orderRepo.GetReturnRequestTablePaginated(offset, limit)
	if err != nil {
		return nil, errors.New("can't fetch return requests")
	}
	if ReturnedRequestTable.ID == 0 {
		return nil, errors.New("currently no return requests")
	}
	return ReturnedRequestTable, nil
}

func (ou *OrderUsecase) ExecuteSuccessfullOrdersPaginated(offset int, limit int) ([]entity.Order, error) {
	SuccessfullOrders, err := ou.orderRepo.GetSuccessfullOrdersPaginated(offset, limit)
	if err != nil {
		return nil, errors.New("can't show successfull orders")
	}
	return SuccessfullOrders, nil
}

func (ou *OrderUsecase) ExecuteSalesInParticularDatePaginated(ParsedTime time.Time, offset int, limit int) ([]entity.Order, error) {
	SuccessfullOrdersInThisDate, err := ou.orderRepo.GetSuccessfullOrderByDate(ParsedTime, offset, limit)
	if err != nil {
		return nil, errors.New("can't fetch sales in this date")
	}
	return SuccessfullOrdersInThisDate, nil
}

func (ou *OrderUsecase) ExecuteTotalOfOrders(userID int) (int, error) {
	TotalOrders, err := ou.orderRepo.GetTotalOrdersOfUser(userID)
	if err != nil {
		return 0, errors.New("can't fetch total of orders")
	}
	return TotalOrders, nil
}

func (ou *OrderUsecase) ExecuteTotalOfOrderedItems(userID int) (int, error) {
	TotalOrderedItems, err := ou.orderRepo.GetTotalOrderedItemsOfUser(userID)
	if err != nil {
		return 0, errors.New("can't fetch total of ordered items")
	}
	return TotalOrderedItems, nil
}

func (ou *OrderUsecase) ExecuteTotalOfCoupons() (int, error) {
	TotalOfAvaialableCoupons, err := ou.orderRepo.GetTotalOfCoupons()
	if err != nil {
		return 0, errors.New("can't fetch total of the available coupons")
	}
	return TotalOfAvaialableCoupons, nil
}

func (ou *OrderUsecase) ExecuteTotalOfSales() (int, error) {
	TotalSales, err := ou.orderRepo.GetTotalOfSales()
	if err != nil {
		return 0, errors.New("can't fetch the total count of the sales")
	}
	return TotalSales, nil
}

func (ou *OrderUsecase) ExecuteTotalOfSalesByDate(Date time.Time) (int, error) {
	TotalSalesByDate, err := ou.orderRepo.GetTotalOfSalesByDate(Date)
	if err != nil {
		return 0, errors.New("can't fetch the total count of sales in this date")
	}
	return TotalSalesByDate, nil
}

func (ou *OrderUsecase) ExecuteTotalOfAllOrders() (int, error) {
	TotalOfAllOrders, err := ou.orderRepo.GetTotalOfAllOrders()
	if err != nil {
		return 0, errors.New("can't fetch the total count of sales in this date")
	}
	return TotalOfAllOrders, nil
}

func (ou *OrderUsecase) ExecuteTotalOfReturnRequests() (int, error) {
	TotalOfRequests, err := ou.orderRepo.GetTotalOfReturnRequests()
	if err != nil {
		return 0, errors.New("can't fetch the total of return requets")
	}
	return TotalOfRequests, nil
}
