//go:build wireinject
// +build wireinject

package di

import (
	"MAXPUMP1/pkg/api/handlers"
	"MAXPUMP1/pkg/db"
	"MAXPUMP1/pkg/repository"
	"MAXPUMP1/pkg/usecase"

	"github.com/google/wire"
)

func InitializeUserApi() *handlers.UserHandler {
	wire.Build(
		db.ConnectDB,
		repository.NewUserRepository,
		repository.NewCategoryRepository,
		repository.NewProductRepository,
		repository.NewCartRepository,
		repository.NewOrderRepository,
		usecase.NewUser,
		usecase.NewCategory,
		usecase.NewProduct,
		usecase.NewCart,
		usecase.NewOrder,
		handlers.NewUserHandler,
	)
	return &handlers.UserHandler{}
}

func InitializeAdminApi() *handlers.AdminHandler {
	wire.Build(
		db.ConnectDB,
		repository.NewAdminRepository,
		repository.NewCategoryRepository,
		repository.NewProductRepository,
		repository.NewOrderRepository,
		usecase.NewAdmin,
		usecase.NewCategory,
		usecase.NewProduct,
		usecase.NewOrder,
		handlers.NewAdminHandler,
	)

	return &handlers.AdminHandler{}
}
