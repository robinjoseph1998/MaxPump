package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
	"MAXPUMP1/pkg/utils"
)

type ProductUsecaseInterface interface {
	ProductCreater(pro entity.Product) (*entity.Product, *entity.Category, error)
	ExecutePaginatedProducts(offset int, limit int) ([]entity.Product, error)
	GetProductByID(ID uint) (*entity.Product, *entity.Category, error)
	ProductUpdate(product *utils.EditProductRequest) (*entity.Product, *entity.Category, error)
	DeleteProduct(ID uint) error
	FetchProductByBrandName(BrandName string, offset int, limit int) ([]entity.Product, error)
	FetchBrandName(BrandName string) (*entity.Product, error)
	FetchPaginatedProductByItemName(ItemName string, offset int, limit int) ([]entity.Product, error)
	ExecuteTotalProducts() (int, error)
	ExecuteTotalOfProductsByBrand(BrandName string) (int, error)
	ExecuteTotalOfProductsByItemName(ItemName string) (int, error)
	FetchItemName(ItemName string) (*entity.Product, error)
}
