package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/model"
)

type UserUsecaseInterface interface {
	ExecuteSignup(use entity.User) (*entity.User, error)
	ExecuteSignupWithOtp(user model.Signup) (string, error)
	ExecuteSignupOtpValidation(key string, otp string) error
	ExecuteLoginWithPassword(phone string, password string) (uint, error)
	FetchProfile(userID int) (*entity.User, []entity.Address, error)
	ExecuteWallet(userID int) (*entity.Wallet, error)
	ExecuteCreateInvoice(userID int) error
}
