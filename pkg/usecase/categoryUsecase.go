package usecase

import (
	"MAXPUMP1/pkg/domain/entity"
	repo "MAXPUMP1/pkg/repository/interfaces"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"errors"
)

type CategoryUsecase struct {
	categoryRepo repo.CategoryInterface
}

func NewCategory(categoryRepo repo.CategoryInterface) use.CategoryUsecaseInterface {
	return &CategoryUsecase{categoryRepo: categoryRepo}

}

func (au *CategoryUsecase) CategoryCreater(cat entity.Category) (*entity.Category, error) {
	name, err := au.categoryRepo.SearchByName(cat.Name)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if name.Name != "" {
		return nil, errors.New("category with this name already exists")
	}
	newCategory := &entity.Category{
		Name:        cat.Name,
		Description: cat.Description,
	}
	err1 := au.categoryRepo.CreateCat(newCategory)

	if err1 != nil {
		return nil, err1
	}
	return newCategory, nil
}

func (au *CategoryUsecase) GetAllCategories() ([]entity.Category, []map[int]int, []map[int]int, error) {
	categories, err := au.categoryRepo.GetAllCategories()
	if err != nil {
		return nil, nil, nil, err
	}
	totalProductswithCategoryID := make([]map[int]int, len(categories))
	totalBrandswithCategoryID := make([]map[int]int, len(categories))
	for index, category := range categories {
		TotalOfProductsPerCategoryId, err := au.categoryRepo.GetTotalProductsInParticularCategory(int(category.ID))
		if err != nil {
			return nil, nil, nil, err
		}
		TotalOfBrandsPerCategoryID, err := au.categoryRepo.GetTotalOfBrandsPerCategory(int(category.ID))
		if err != nil {
			return nil, nil, nil, err
		}
		categoryTotalProducts := map[int]int{int(category.ID): TotalOfProductsPerCategoryId}
		totalProductswithCategoryID[index] = categoryTotalProducts

		categoryTotalBrands := map[int]int{int(category.ID): TotalOfBrandsPerCategoryID}
		totalBrandswithCategoryID[index] = categoryTotalBrands
	}
	return categories, totalProductswithCategoryID, totalBrandswithCategoryID, nil
}

func (au *CategoryUsecase) GetCategoryByID(ID uint) (*entity.Category, error) {
	category, err := au.categoryRepo.GetCategoryByID(ID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (au *CategoryUsecase) CategoryUpdate(UpdatedCategory *entity.Category) (*entity.Category, error) {
	updatedcategory, err := au.categoryRepo.UpdateCategory(UpdatedCategory)
	if err != nil {
		return nil, err
	}
	return updatedcategory, nil
}

func (au *CategoryUsecase) CategoryWithPaginatedProducts(Name string, offset int, limit int) ([]entity.Product, *entity.Category, error) {
	category, err := au.categoryRepo.SearchByName(Name)
	if err != nil {
		return nil, nil, err
	}
	if category.Name == "" {
		return nil, nil, errors.New("category with this name not exist")
	}
	CategoryAndProducts, err := au.categoryRepo.GetCategoryWithPaginatedProducts(category.ID, offset, limit)
	if err != nil {
		return nil, nil, err
	}
	return CategoryAndProducts, category, nil
}

func (au *CategoryUsecase) DeleteCategory(ID uint) error {
	err := au.categoryRepo.CategoryDelete(ID)
	if err != nil {
		return err
	}
	return nil
}

func (au *CategoryUsecase) FetchCategoryByName(Name string) (*entity.Category, int64, error) {
	FetchedCategory, err := au.categoryRepo.SearchByName(Name)
	if err != nil {
		return nil, 0, err
	}
	if FetchedCategory.Name == "" {
		return nil, 0, errors.New("category with this name not exist")
	}
	categoryID := FetchedCategory.ID
	var ProductCount int64
	ProductCount, err = au.categoryRepo.CountOfProducts(categoryID)
	if err != nil {
		return nil, 0, errors.New("can't fetch the products count")
	}
	return FetchedCategory, ProductCount, nil
}

func (au *CategoryUsecase) ExecutePaginatedCategories(offset int, limit int) ([]entity.Category, error) {
	paginatedCategories, err := au.categoryRepo.GetPaginatedCategories(offset, limit)
	if err != nil {
		return nil, errors.New("can't execute paginated categories")
	}
	return paginatedCategories, nil
}

func (au *CategoryUsecase) ExecuteTotalCategories() (int, error) {
	TotalOfCategories, err := au.categoryRepo.GetTotalOfCategories()
	if err != nil {
		return 0, errors.New("can't get the total of categories")
	}
	return TotalOfCategories, nil
}

func (au *CategoryUsecase) ExecuteTotalProductsInTheParticularCategory(CategoryId int) (int, error) {
	TotalOfProductsInThisCategory, err := au.categoryRepo.GetTotalProductsInParticularCategory(CategoryId)
	if err != nil {
		return 0, errors.New("can't get the total of products in this category")
	}
	return TotalOfProductsInThisCategory, nil
}
