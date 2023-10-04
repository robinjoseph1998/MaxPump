package interfaces

import "MAXPUMP1/pkg/domain/entity"

type ProductInterface interface {
	SearchByBrandNameAndItem(BrandName string, item string) (*entity.Product, error)
	CreateProduct(product *entity.Product) (*entity.Product, error)
	GetCategoryByID(id int) (*entity.Category, error)
	GetPaginatedProducts(offset int, limit int) ([]entity.Product, error)
	GetProductByID(id uint) (*entity.Product, error)
	ProductUpdate(product *entity.Product) (*entity.Product, error)
	ProductDelete(id uint) error
	GetPaginatedProductsByBrandName(BrandName string, offset int, limit int) ([]entity.Product, error)
	GetByBrand(BrandName string) (*entity.Product, error)
	GetPaginatedProductsByItemName(ItemName string, offset int, limit int) ([]entity.Product, error)
	GetTotalOfProducts() (int, error)
	GetTotalOfProductsByBrand(BrandName string) (int, error)
	GetTotalOfProductsByItemName(ItemName string) (int, error)
	GetByItem(ItemName string) (*entity.Product, error)
}
