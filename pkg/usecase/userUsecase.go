package usecase

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/model"
	repo "MAXPUMP1/pkg/repository/interfaces"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"MAXPUMP1/pkg/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo repo.UserInterface
}

func NewUser(userRepo repo.UserInterface) use.UserUsecaseInterface {
	return &UserUsecase{userRepo: userRepo}
}

func (uu *UserUsecase) ExecuteSignup(user entity.User) (*entity.User, error) {
	email, err := uu.userRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if email.Email != "" {
		return nil, errors.New("user with this email already exists")
	}
	phone, err := uu.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if phone.Phone != "" {
		return nil, errors.New("user with this phone no already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := &entity.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  string(hashedPassword),
	}
	err1 := uu.userRepo.Create(newUser)
	if err1 != nil {
		return nil, err1
	}
	return newUser, nil
}

func (uu *UserUsecase) ExecuteSignupWithOtp(user model.Signup) (string, error) {
	var otpKey entity.OtpKey
	email, err := uu.userRepo.GetByEmail(user.Email)
	if err != nil {
		return "", errors.New("error with server")
	}
	if email.Email != "" {
		return "", errors.New("user with this email already exists")
	}
	phone, err := uu.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return "", errors.New("error with server")
	}
	if phone.Phone != "" {
		return "", errors.New("user with this phone no already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)
	key, err := utils.SendOtp(user.Phone)
	if err != nil {
		return "", err
	} else {
		err := uu.userRepo.CreateSignup(&user)
		if err != nil {
			return "", errors.New("error creating the signup")
		}
		otpKey.Key = key
		otpKey.Phone = user.Phone
		err = uu.userRepo.CreateOtpKey(&otpKey)
		if err != nil {
			return "", err
		}
		return key, nil
	}
}

func (uu *UserUsecase) ExecuteSignupOtpValidation(key string, otp string) error {
	result, err := uu.userRepo.GetByKey(key)
	if err != nil {
		return err
	}
	user, err := uu.userRepo.GetSignupByPhone(result.Phone)
	if err != nil {
		return err
	}
	err = utils.CheckOtp(result.Phone, otp)
	if err != nil {
		return err
	} else {
		newUser := &entity.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			Password:  user.Password,
		}
		err1 := uu.userRepo.Create(newUser)
		if err1 != nil {
			return err1
		} else {
			return nil
		}
	}
}

func (uu *UserUsecase) ExecuteLoginWithPassword(phone string, password string) (uint, error) {
	user, err := uu.userRepo.GetByPhone(phone)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("user with this phone number not exist")
	}
	blocked := uu.userRepo.IsBlocked(phone)
	if blocked {
		return 0, errors.New("your account is blocked")
	}
	permission, err := uu.userRepo.CheckPermission(user)
	if err != nil {
		return 0, errors.New("error in checking permission" + err.Error())
	}
	if !permission {
		return 0, errors.New("user permission denied")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("invalid password")
	} else {
		return user.ID, nil
	}
}

func (uu *UserUsecase) FetchProfile(userID int) (*entity.User, []entity.Address, error) {
	user, err := uu.userRepo.GetByID(userID)
	if err != nil {
		return nil, nil, errors.New("error in get user")
	}
	address, err := uu.userRepo.GetAddressesByUserID(userID)
	if err != nil {
		return nil, nil, errors.New("error in get address")
	}
	return user, address, nil
}

func (uu *UserUsecase) ExecuteWallet(userID int) (*entity.Wallet, error) {
	MyWallet, err := uu.userRepo.GetMyWallet(userID)
	if err != nil {
		return nil, errors.New("can't get wallet")
	}
	return MyWallet, nil
}

func (uu *UserUsecase) ExecuteCreateInvoice(userID int) error {
	Order, err := uu.userRepo.GetOrdersByUserID(userID)
	if err != nil {
		return errors.New("can't fetch order")
	}
	BillingAddress, err := uu.userRepo.GetAddressByOrder(Order.AddressId)
	if err != nil {
		return errors.New("can't get address")
	}
	User, err := uu.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("can't get user")
	}
	err = utils.InvoiceGenerator(Order, BillingAddress, User)
	if err != nil {
		return errors.New("can't generate invoice")
	}
	return nil
}
