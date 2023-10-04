package interfaces

import "MAXPUMP1/pkg/domain/entity"

type CategoryInterface interface {
	SearchByName(name string) (*entity.Category, error)
	CreateCat(category *entity.Category) error
	GetAllCategories() ([]entity.Category, error)
	GetCategoryByID(id uint) (*entity.Category, error)
	UpdateCategory(category *entity.Category) (*entity.Category, error)
	GetCategoryWithPaginatedProducts(ID uint, offset int, limit int) ([]entity.Product, error)
	CategoryDelete(ID uint) error
	CountOfProducts(categoryID uint) (int64, error)
	GetPaginatedCategories(offset int, limit int) ([]entity.Category, error)
	GetTotalOfCategories() (int, error)
	GetTotalProductsInParticularCategory(ID int) (int, error)
	GetTotalOfBrandsPerCategory(ID int) (int, error)
}
