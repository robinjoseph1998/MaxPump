package repository

import (
	"MAXPUMP1/pkg/domain/entity"
	"errors"

	repo "MAXPUMP1/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repo.ProductInterface {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) CreateProduct(product *entity.Product) (*entity.Product, error) {
	var Product entity.Product
	err := pr.db.Create(product).Scan(&Product).Error
	if err != nil {
		return nil, err
	}
	return &Product, nil
}

func (pr *ProductRepository) GetCategoryByID(id int) (*entity.Category, error) {
	var result entity.Category
	err := pr.db.Raw("SELECT * FROM categories WHERE id=?", id).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (pr *ProductRepository) SearchByBrandNameAndItem(BrandName string, item string) (*entity.Product, error) {
	var result entity.Product
	err := pr.db.Raw("SELECT * FROM products WHERE brand_name=? AND item=?", BrandName, item).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, err
}

func (pr *ProductRepository) GetPaginatedProducts(offset int, limit int) ([]entity.Product, error) {
	var products []entity.Product
	err := pr.db.Raw("SELECT * FROM products OFFSET ? LIMIT ?", offset, limit).Scan(&products).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return products, nil
}

func (pr *ProductRepository) GetProductByID(id uint) (*entity.Product, error) {
	var result entity.Product
	err := pr.db.Raw("SELECT * FROM products WHERE id=?", id).Scan(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (pr *ProductRepository) ProductUpdate(product *entity.Product) (*entity.Product, error) {
	var updatedpro entity.Product
	err := pr.db.Exec("UPDATE products SET brand_name=?,description=?,price=?,quantity=?,category_id=? WHERE id = ?", product.Brand_Name, product.Description, product.Price, product.Quantity, product.CategoryID, product.ID).Error
	if err != nil {
		return nil, err
	}
	err = pr.db.Raw("SELECT * FROM products WHERE id=?", product.ID).Scan(&updatedpro).Error
	if err != nil {
		return nil, err
	}
	return &updatedpro, nil
}

func (pr *ProductRepository) ProductDelete(id uint) error {
	err := pr.db.Exec("DELETE FROM products WHERE id=?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetPaginatedProductsByBrandName(BrandName string, offset int, limit int) ([]entity.Product, error) {
	var FetchedProducts []entity.Product
	err := pr.db.Raw("SELECT * FROM products WHERE brand_name=? OFFSET ? LIMIT ?", BrandName, offset, limit).Scan(&FetchedProducts).Error
	if err != nil {
		return nil, err
	}
	return FetchedProducts, nil
}

func (pr *ProductRepository) GetByBrand(BrandName string) (*entity.Product, error) {
	var FetchedBrand entity.Product
	err := pr.db.Raw("SELECT * FROM products WHERE brand_name=?", BrandName).Scan(&FetchedBrand).Error
	if err != nil {
		return nil, err
	}
	return &FetchedBrand, nil
}

func (pr *ProductRepository) GetByItem(ItemName string) (*entity.Product, error) {
	var FetchedItem entity.Product
	err := pr.db.Raw("SELECT * FROM products WHERE item=?", ItemName).Scan(&FetchedItem).Error
	if err != nil {
		return nil, err
	}
	return &FetchedItem, nil
}

func (pr *ProductRepository) GetPaginatedProductsByItemName(ItemName string, offset int, limit int) ([]entity.Product, error) {
	var FetchedProducts []entity.Product
	err := pr.db.Raw("SELECT * FROM products WHERE item=? OFFSET ? LIMIT ?", ItemName, offset, limit).Scan(&FetchedProducts).Error
	if err != nil {
		return nil, err
	}
	return FetchedProducts, nil
}

func (pr *ProductRepository) GetTotalOfProducts() (int, error) {
	var Count int
	err := pr.db.Raw("SELECT COUNT(*) FROM products").Scan(&Count).Error
	if err != nil {
		return 0, err
	}
	return Count, nil
}

func (pr *ProductRepository) GetTotalOfProductsByBrand(BrandName string) (int, error) {
	var count int
	err := pr.db.Raw("SELECT COUNT(*) FROM products WHERE brand_name=?", BrandName).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pr *ProductRepository) GetTotalOfProductsByItemName(ItemName string) (int, error) {
	var count int
	err := pr.db.Raw("SELECT COUNT(*) FROM products WHERE item=?", ItemName).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
