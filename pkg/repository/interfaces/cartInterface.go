package interfaces

import "MAXPUMP1/pkg/domain/entity"

type CartInterface interface {
	CheckQuantity(productID int) (int, error)
	GetProductByID(ID uint) (*entity.Product, error)
	GetCartByUserID(ID int) (*entity.Cart, bool, error)
	CreateCart(userID int) (*entity.Cart, error)
	UpdateCartPriceAndCount(cart *entity.Cart) (*entity.Cart, error)
	CompareProduct(ProductID int, cartID int, UserID int) (bool, error)
	GetCartItemByProductID(ProductID int, cartID int, UserID int) (*entity.CartItem, error)
	UpdateQuantity(UserID int, ProductID int, qty int) error
	CreateCartItem(cartItem *entity.CartItem) (*entity.CartItem, error)
	UpdateCartItem(cartItem *entity.CartItem) (*entity.CartItem, error)
	UpdateCart(cart *entity.Cart) (*entity.Cart, error)
	GetAllCarts(userID int) ([]entity.Cart, error)
	GetPaginatedCartItems(userID int, offset int, limit int) ([]entity.CartItem, error)
	DeleteProductFromCart(productID int, UserID int) error
	CreateAddress(UserAddress *entity.Address) (*entity.Address, error)
	GetAddressByUserId(userID int) (*entity.Address, error)
	EditAddress(ExistingAddress *entity.Address, userID int) (*entity.Address, error)
	GetAddressesByUserID(userID int) ([]entity.Address, error)
	GetTotalOfCartItems(userID int) (int, error)
}
