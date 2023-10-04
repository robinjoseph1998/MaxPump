package usecase

import (
	"MAXPUMP1/pkg/domain/entity"
	repo "MAXPUMP1/pkg/repository/interfaces"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"MAXPUMP1/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

type CartUsecase struct {
	cartRepo repo.CartInterface
}

func NewCart(cartRepo repo.CartInterface) use.CartUsecaseInterface {

	return &CartUsecase{cartRepo: cartRepo}

}

func (cu *CartUsecase) CheckQuantity(productID int) (int, error) {
	quantity, err := cu.cartRepo.CheckQuantity(productID)
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

func (cu *CartUsecase) GetProductByID(id uint) (*entity.Product, error) {
	product, err := cu.cartRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (cu *CartUsecase) AddToCart(UserID int, ProductID int, qty int) (*entity.Cart, *entity.CartItem, error) {
	var cart *entity.Cart
	CartExist := true
	//Fetching User Cart By UserID
	cart, CartExist, err := cu.cartRepo.GetCartByUserID(UserID)
	if err != nil {
		return nil, nil, err
	}
	if !CartExist {
		var Newcart *entity.Cart
		//if there is no Existing Cart Creating New Cart
		Newcart, err = cu.cartRepo.CreateCart(UserID)
		if err != nil {
			return nil, nil, errors.New("failed To create cart")
		}
		cart = Newcart
		//Fetching Product By Product ID
		product, err := cu.cartRepo.GetProductByID(uint(ProductID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, errors.New("there is no product in this id")
			}
			return nil, nil, err
		}
		cartItem := &entity.CartItem{
			CartId:        int(cart.ID),
			UserID:        UserID,
			ProductItemID: int(product.ID),
			Qty:           qty,
			ProductName:   product.Brand_Name,
			Price:         product.Price,
		}
		CreatedCartItem, err := cu.cartRepo.CreateCartItem(cartItem)
		if err != nil {
			return nil, nil, errors.New("can't create cart item")
		}
		//Creating a New Cart and CartItem Based on parameters
		cart.TotalProducts += 1
		cart.UserID = UserID
		cart.TotalPrice += float64(CreatedCartItem.Qty) * product.Price
		cart, err = cu.cartRepo.UpdateCart(cart)
		if err != nil {
			return nil, nil, errors.New("can't update count of products added in cart")
		}
		return cart, CreatedCartItem, nil
	}
	if CartExist {
		//Fetching Product By Product ID
		product, err := cu.cartRepo.GetProductByID(uint(ProductID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, errors.New("there is no product in this id")
			}
			return nil, nil, err
		}
		//Comparing the product is already existed in the user Cart
		ProductExists, err := cu.cartRepo.CompareProduct(ProductID, int(cart.ID), UserID)
		if err != nil {
			return nil, nil, err
		}
		if ProductExists {
			// If the product already exists
			cartItem, err := cu.cartRepo.GetCartItemByProductID(ProductID, int(cart.ID), UserID)
			if err != nil {
				return nil, nil, err
			}
			// Update the quantity
			cartItem.Qty += qty
			// Update the cartItem in the database
			updatedCartItem, err := cu.cartRepo.UpdateCartItem(cartItem)
			if err != nil {
				return nil, nil, err
			}
			// Updating the total price of the cart in the database
			cart.TotalPrice += float64(updatedCartItem.Qty) * product.Price
			updatedCart, err := cu.cartRepo.UpdateCartPriceAndCount(cart)
			if err != nil {
				return nil, nil, err
			}
			return updatedCart, updatedCartItem, nil
		}
		if !ProductExists {
			cartItem := &entity.CartItem{
				CartId:        int(cart.ID),
				UserID:        UserID,
				ProductItemID: int(product.ID),
				Qty:           qty,
				ProductName:   product.Brand_Name,
				Price:         product.Price,
			}
			CreatedCartItem, err := cu.cartRepo.CreateCartItem(cartItem)
			if err != nil {
				return nil, nil, errors.New("can't create cart item")
			}
			// Updating the total price  and product count of the cart in the database
			cart.TotalProducts++
			cart.TotalPrice += float64(CreatedCartItem.Qty) * product.Price
			NewCart, err := cu.cartRepo.UpdateCartPriceAndCount(cart)
			if err != nil {
				return nil, nil, err
			}
			return NewCart, CreatedCartItem, nil
		}
	}
	return nil, nil, errors.New("can't create or update cart")
}

func (cu *CartUsecase) ExecuteCart(userID int) ([]entity.Cart, error) {
	Carts, err := cu.cartRepo.GetAllCarts(userID)
	if err != nil {
		return nil, err
	}
	return Carts, nil

}

func (cu *CartUsecase) ListPaginatedCartItems(userID int, offset int, limit int) ([]entity.CartItem, error) {
	CartItems, err := cu.cartRepo.GetPaginatedCartItems(userID, offset, limit)
	if err != nil {
		return nil, err
	}
	return CartItems, nil
}

func (cu *CartUsecase) DropProduct(productID int, UserId int) error {
	err := cu.cartRepo.DeleteProductFromCart(productID, UserId)
	if err != nil {
		return err
	}
	return nil
}

func (cu *CartUsecase) AddAddress(UserAddress *entity.Address) (*entity.Address, error) {
	AddedAddress, err := cu.cartRepo.CreateAddress(UserAddress)
	if err != nil {
		return nil, err
	}
	return AddedAddress, nil
}

func (cu *CartUsecase) EditAddress(request utils.EditAddressRequest, userID int) (*entity.Address, error) {
	var EditedAddress *entity.Address
	ExistingAddress, err := cu.cartRepo.GetAddressByUserId(userID)
	if err != nil {
		return nil, errors.New("can't fetch Address in Database")
	}
	if request.HouseName != nil {
		ExistingAddress.HouseName = *request.HouseName
	}
	if request.Street != nil {
		ExistingAddress.Street = *request.Street
	}
	if request.City != nil {
		ExistingAddress.City = *request.City
	}
	if request.District != nil {
		ExistingAddress.District = *request.District
	}
	if request.State != nil {
		ExistingAddress.State = *request.State
	}
	if request.Pincode != nil {
		ExistingAddress.Pincode = *request.Pincode
	}
	if request.Landmark != nil {
		ExistingAddress.Landmark = *request.Landmark
	}
	if EditedAddress, err = cu.cartRepo.EditAddress(ExistingAddress, userID); err != nil {
		return nil, errors.New("can't edit the address")
	}
	return EditedAddress, nil
}

func (cu *CartUsecase) UpdateAddress(request *entity.Address, userID int) (*entity.Address, error) {
	updatedAddress, err := cu.cartRepo.EditAddress(request, userID)
	if err != nil {
		return nil, errors.New("can't perform updation")
	}
	return updatedAddress, nil
}

func (cu *CartUsecase) GetAllAddresses(userID int) ([]entity.Address, error) {
	AllAddress, err := cu.cartRepo.GetAddressesByUserID(userID)
	if err != nil {
		return nil, errors.New("can't fetch addresses")
	}
	return AllAddress, nil
}

func (cu *CartUsecase) ExecuteTotalItemsInCart(userID int) (int, error) {
	TotalItemsIncart, err := cu.cartRepo.GetTotalOfCartItems(userID)
	if err != nil {
		return 0, errors.New("can't get total of cart items")
	}
	return TotalItemsIncart, nil
}
