package repository

import (
	"MAXPUMP1/pkg/domain/entity"
	"errors"

	repo "MAXPUMP1/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repo.CategoryInterface {
	return &CategoryRepository{db: db}

}

func (cr *CategoryRepository) CreateCat(category *entity.Category) error {
	return cr.db.Create(category).Error

}
func (cr *CategoryRepository) SearchByName(name string) (*entity.Category, error) {
	var result entity.Category
	err := cr.db.Raw("SELECT * FROM categories WHERE name=?", name).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (cr *CategoryRepository) GetAllCategories() ([]entity.Category, error) {
	var category []entity.Category
	err := cr.db.Raw("SELECT * FROM categories").Scan(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func (cr *CategoryRepository) GetCategoryByID(id uint) (*entity.Category, error) {
	var result entity.Category
	err := cr.db.Raw("SELECT * FROM categories WHERE id=?", id).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (cr *CategoryRepository) UpdateCategory(category *entity.Category) (*entity.Category, error) {
	var updatedCat entity.Category
	err := cr.db.Raw("UPDATE categories SET name=?,description=? WHERE id = ?", category.Name, category.Description, category.ID).Error
	if err != nil {
		return nil, err
	}
	err = cr.db.Raw("SELECT * FROM categories WHERE id=?", category.ID).Scan(&updatedCat).Error
	if err != nil {
		return nil, err
	}
	return &updatedCat, nil
}

func (cr *CategoryRepository) GetCategoryWithPaginatedProducts(ID uint, offset int, limit int) ([]entity.Product, error) {
	var products []entity.Product
	err := cr.db.Raw("SELECT * FROM products WHERE category_id = ? OFFSET ? LIMIT ?", ID, offset, limit).Scan(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (cr *CategoryRepository) CategoryDelete(ID uint) error {
	err := cr.db.Exec("DELETE FROM categories WHERE id=?", ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr *CategoryRepository) CountOfProducts(categoryID uint) (int64, error) {
	var count int64
	err := cr.db.Raw("SELECT COUNT(*) FROM products WHERE category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (cr *CategoryRepository) GetPaginatedCategories(offset int, limit int) ([]entity.Category, error) {
	var paginatedCategories []entity.Category
	err := cr.db.Raw("SELECT * FROM categories OFFSET ? LIMIT ?", offset, limit).Scan(&paginatedCategories).Error
	if err != nil {
		return nil, err
	}
	return paginatedCategories, nil
}

func (cr *CategoryRepository) GetTotalOfCategories() (int, error) {
	var Count int
	err := cr.db.Raw("SELECT COUNT(*) FROM categories").Scan(&Count).Error
	if err != nil {
		return 0, err
	}
	return Count, nil
}

func (cr *CategoryRepository) GetTotalProductsInParticularCategory(ID int) (int, error) {
	var count int
	err := cr.db.Raw("SELECT COUNT(*) FROM products WHERE category_id=?", ID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (cr *CategoryRepository) GetTotalOfBrandsPerCategory(ID int) (int, error) {
	var count int
	err := cr.db.Raw("SELECT COUNT(brand_name) FROM products WHERE category_id=?", ID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
