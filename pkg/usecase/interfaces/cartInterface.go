package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/utils"
)

type CartUsecaseInterface interface {
	CheckQuantity(productID int) (int, error)
	GetProductByID(id uint) (*entity.Product, error)
	AddToCart(UserID int, ProductID int, qty int) (*entity.Cart, *entity.CartItem, error)
	ExecuteCart(userID int) ([]entity.Cart, error)
	ListPaginatedCartItems(userID int, offset int, limit int) ([]entity.CartItem, error)
	DropProduct(productID int, UserId int) error
	AddAddress(UserAddress *entity.Address) (*entity.Address, error)
	EditAddress(request utils.EditAddressRequest, userID int) (*entity.Address, error)
	GetAllAddresses(userID int) ([]entity.Address, error)
	ExecuteTotalItemsInCart(userID int) (int, error)
}
