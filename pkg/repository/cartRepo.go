package repository

import (
	"MAXPUMP1/pkg/domain/entity"
	repo "MAXPUMP1/pkg/repository/interfaces"
	"errors"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) repo.CartInterface {
	return &CartRepository{db: db}
}

func (cr *CartRepository) CheckQuantity(productID int) (int, error) {
	var value int
	err := cr.db.Raw("SELECT quantity FROM products WHERE id=?", productID).Scan(&value).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return value, err
}

func (cr *CartRepository) GetProductByID(ID uint) (*entity.Product, error) {
	var result entity.Product
	err := cr.db.Raw("SELECT * FROM products WHERE id=?", ID).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (cr *CartRepository) GetCartByUserID(ID int) (*entity.Cart, bool, error) {
	var cart entity.Cart
	err := cr.db.Raw("SELECT * FROM carts WHERE user_id=?", ID).Scan(&cart).Error
	if err != nil {
		return nil, false, err
	}
	if cart.UserID == 0 {
		return nil, false, err
	}
	return &cart, true, nil
}

func (cr *CartRepository) CreateCart(userID int) (*entity.Cart, error) {
	cart := &entity.Cart{
		UserID: userID,
	}
	if err := cr.db.Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (cr *CartRepository) CreateCartItem(cartItem *entity.CartItem) (*entity.CartItem, error) {
	if err := cr.db.Create(cartItem).Error; err != nil {
		return nil, err
	}
	return cartItem, nil
}

func (cr *CartRepository) GetCartItemByProductID(ProductID int, cartID int, UserID int) (*entity.CartItem, error) {
	var result entity.CartItem
	err := cr.db.Raw("SELECT * FROM cart_items WHERE product_item_id=? AND cart_id=? AND user_id=?", ProductID, cartID, UserID).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cr *CartRepository) CompareProduct(ProductID int, cartID int, UserID int) (bool, error) {
	var count int64
	err := cr.db.Raw("SELECT COUNT(*) FROM cart_items WHERE product_item_id = ? AND cart_id = ? AND user_id = ?", ProductID, cartID, UserID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (cr *CartRepository) UpdateCartItem(cartItem *entity.CartItem) (*entity.CartItem, error) {
	var CartItem entity.CartItem
	err := cr.db.Exec("UPDATE cart_items SET qty = ? WHERE id = ?", cartItem.Qty, cartItem.ID).Error
	if err != nil {
		return nil, err
	}
	err1 := cr.db.Raw("SELECT * FROM cart_items WHERE id=?", cartItem.ID).Scan(&CartItem).Error
	if err != nil {
		return nil, err1
	}
	return &CartItem, nil
}

func (cr *CartRepository) UpdateCartPriceAndCount(cart *entity.Cart) (*entity.Cart, error) {
	var CarT entity.Cart
	err := cr.db.Exec("UPDATE carts SET total_price = ?,total_products=? WHERE id = ?", cart.TotalPrice, cart.TotalProducts, cart.ID).Error
	if err != nil {
		return nil, err
	}
	err1 := cr.db.Raw("SELECT * FROM carts WHERE id=?", cart.ID).Scan(&CarT).Error
	if err != nil {
		return nil, err1
	}
	return &CarT, nil
}

func (cr *CartRepository) UpdateCart(cart *entity.Cart) (*entity.Cart, error) {
	var UpdatedCart entity.Cart
	err := cr.db.Exec("UPDATE carts SET total_products = ?,user_id=?,total_price=? WHERE id = ?", cart.TotalProducts, cart.UserID, cart.TotalPrice, cart.ID).Error
	if err != nil {
		return nil, err
	}
	err1 := cr.db.Raw("SELECT * FROM carts WHERE id=?", cart.ID).Scan(&UpdatedCart).Error
	if err1 != nil {
		return nil, err1
	}
	return &UpdatedCart, nil
}

func (cr *CartRepository) UpdateQuantity(UserID int, ProductID int, qty int) error {
	err := cr.db.Exec("UPDATE carts SET qty = ? WHERE product_item_id=? AND user_id=?", qty, ProductID, UserID).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr *CartRepository) GetAllCarts(UserId int) ([]entity.Cart, error) {
	var ListCart []entity.Cart
	err := cr.db.Raw("SELECT * FROM carts WHERE user_id=?", UserId).Scan(&ListCart).Error
	if err != nil {
		return nil, err
	}
	return ListCart, nil
}

func (cr *CartRepository) GetPaginatedCartItems(userID int, offset int, limit int) ([]entity.CartItem, error) {
	var ListItems []entity.CartItem
	err := cr.db.Raw("SELECT * FROM cart_items WHERE user_id=? OFFSET ? LIMIT ?", userID, offset, limit).Scan(&ListItems).Error
	if err != nil {
		return nil, err
	}
	return ListItems, nil
}

func (cr *CartRepository) DeleteProductFromCart(productID int, UserID int) error {
	err := cr.db.Exec("DELETE FROM cart_items WHERE product_item_id = ? AND user_id = ? ", productID, UserID).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr *CartRepository) CreateAddress(UserAddress *entity.Address) (*entity.Address, error) {
	var AddedAddress entity.Address
	err := cr.db.Create(UserAddress).Scan(&AddedAddress).Error
	if err != nil {
		return nil, err
	}
	return &AddedAddress, nil
}

func (cr *CartRepository) GetAddressByUserId(userID int) (*entity.Address, error) {
	var address entity.Address
	err := cr.db.Raw("SELECT * FROM addresses WHERE user_id=?", userID).Scan(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &address, nil
}

func (cr *CartRepository) EditAddress(ExistingAddress *entity.Address, UserID int) (*entity.Address, error) {
	var address entity.Address
	err := cr.db.Raw("UPDATE addresses SET house_name=?,street=?,city=?,district=?,state=?,pincode=?,landmark=? WHERE user_id=?", ExistingAddress.HouseName,
		ExistingAddress.Street,
		ExistingAddress.City,
		ExistingAddress.District,
		ExistingAddress.State,
		ExistingAddress.Pincode,
		ExistingAddress.Landmark, UserID).Error
	if err != nil {
		return nil, err
	}
	err = cr.db.Raw("SELECT * FROM addresses WHERE user_id=?", UserID).Scan(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (cr *CartRepository) GetAddressesByUserID(userID int) ([]entity.Address, error) {
	var AllAddress []entity.Address
	err := cr.db.Raw("SELECT * FROM addresses WHERE user_id=?", userID).Scan(&AllAddress).Error
	if err != nil {
		return nil, err
	}
	return AllAddress, nil
}

func (cr *CartRepository) GetTotalOfCartItems(userID int) (int, error) {
	var count int
	err := cr.db.Raw("SELECT COUNT(*) FROM cart_items WHERE user_id=?", userID).Scan(&count).Error
	if err != nil {
		return 0, nil
	}
	return count, nil
}
