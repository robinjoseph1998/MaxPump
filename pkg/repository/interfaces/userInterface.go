package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/model"
)

type UserInterface interface {
	GetByID(id int) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByPhone(Phone string) (*entity.User, error)
	CreateSignup(user *model.Signup) error
	GetSignupByPhone(phone string) (*model.Signup, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error
	CreateOtpKey(otpKey *entity.OtpKey) error
	GetByKey(key string) (*entity.OtpKey, error)
	CheckPermission(user *entity.User) (bool, error)
	IsBlocked(phone string) bool
	GetAddressesByUserID(userID int) ([]entity.Address, error)
	GetMyWallet(userID int) (*entity.Wallet, error)
	GetOrdersByUserID(userID int) (*entity.Order, error)
	GetAddressByOrder(AddressId int) (*entity.Address, error)
}
