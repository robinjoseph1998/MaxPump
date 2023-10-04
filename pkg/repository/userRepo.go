package repository

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/model"
	"errors"
	"fmt"

	repo "MAXPUMP1/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repo.UserInterface {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetByID(id int) (*entity.User, error) {
	var user entity.User
	err := ur.db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	// Execute the raw query
	err := ur.db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetByPhone(Phone string) (*entity.User, error) {
	var user entity.User
	err := ur.db.Raw("SELECT * FROM users WHERE phone = ?", Phone).Scan(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) CheckPermission(user *entity.User) (bool, error) {
	result := ur.db.Raw("SELECT * FROM users WHERE phone = ?", user.Phone).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	permission := user.Permission
	return permission, nil
}

func (ur *UserRepository) GetByKey(key string) (*entity.OtpKey, error) {
	var otpKey entity.OtpKey
	result := ur.db.Raw("SELECT * FROM otp_keys WHERE key = ?", key).First(&otpKey)
	fmt.Println("otpkey from userRepo", otpKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &otpKey, nil
}

func (ur *UserRepository) CreateSignup(user *model.Signup) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) GetSignupByPhone(phone string) (*model.Signup, error) {
	var user model.Signup
	result := ur.db.Raw("SELECT * FROM signups WHERE phone = ?", phone).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetAddressesByUserID(userID int) ([]entity.Address, error) {
	var address []entity.Address
	err := ur.db.Raw("SELECT * FROM addresses WHERE user_id=?", userID).Scan(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return address, nil
}

func (ur *UserRepository) IsBlocked(phone string) bool {
	var user entity.User
	result := ur.db.Raw("SELECT * FROM users WHERE phone = ? AND blocked = ?", phone, true).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}
	return user.Blocked
}

func (ur *UserRepository) Create(user *entity.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) Update(user *entity.User) error {
	return ur.db.Updates(user).Error
}

func (ur *UserRepository) Delete(id uint) error {
	return ur.db.Delete(&entity.User{}, id).Error
}

func (ar *UserRepository) CreateOtpKey(otpKey *entity.OtpKey) error {
	return ar.db.Create(otpKey).Error
}

func (ur *UserRepository) GetAllCoupons() ([]entity.Coupon, error) {
	var Coupons []entity.Coupon
	err := ur.db.Raw("SELECT * FROM coupons").Scan(&Coupons).Error
	if err != nil {
		return nil, err
	}
	return Coupons, nil
}

func (ur *UserRepository) GetMyWallet(userID int) (*entity.Wallet, error) {
	var Wallet entity.Wallet
	err := ur.db.Raw("SELECT * from wallets WHERE user_id=?", userID).Scan(&Wallet).Error
	if err != nil {
		return nil, err
	}
	return &Wallet, nil
}

func (ur *UserRepository) GetOrdersByUserID(userID int) (*entity.Order, error) {
	var Order entity.Order
	err := ur.db.Raw("SELECT * FROM orders WHERE user_id=?", userID).Scan(&Order).Error
	if err != nil {
		return nil, err
	}
	return &Order, nil
}

func (ur *UserRepository) GetAddressByOrder(AddressId int) (*entity.Address, error) {
	var Address entity.Address
	err := ur.db.Raw("SELECT * FROM addresses WHERE id=?", AddressId).Scan(&Address).Error
	if err != nil {
		return nil, err
	}
	return &Address, nil
}
