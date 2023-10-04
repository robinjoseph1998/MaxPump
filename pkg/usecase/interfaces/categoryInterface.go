package interfaces

import (
	"MAXPUMP1/pkg/domain/entity"
)

type CategoryUsecaseInterface interface {
	CategoryCreater(cat entity.Category) (*entity.Category, error)
	GetAllCategories() ([]entity.Category, []map[int]int, []map[int]int, error)
	GetCategoryByID(ID uint) (*entity.Category, error)
	CategoryUpdate(category *entity.Category) (*entity.Category, error)
	CategoryWithPaginatedProducts(Name string, offset int, limit int) ([]entity.Product, *entity.Category, error)
	DeleteCategory(ID uint) error
	FetchCategoryByName(Name string) (*entity.Category, int64, error)
	ExecutePaginatedCategories(offset int, limit int) ([]entity.Category, error)
	ExecuteTotalCategories() (int, error)
	ExecuteTotalProductsInTheParticularCategory(CategoryId int) (int, error)
}
