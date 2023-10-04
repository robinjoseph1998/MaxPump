package usecase

import (
	"MAXPUMP1/pkg/domain/entity"
	repo "MAXPUMP1/pkg/repository/interfaces"
	use "MAXPUMP1/pkg/usecase/interfaces"
	"MAXPUMP1/pkg/utils"
	"errors"
)

type ProductUsecase struct {
	productRepo repo.ProductInterface
}

func NewProduct(productRepo repo.ProductInterface) use.ProductUsecaseInterface {
	return &ProductUsecase{productRepo: productRepo}
}

func (pu *ProductUsecase) ProductCreater(pro entity.Product) (*entity.Product, *entity.Category, error) {
	name, err := pu.productRepo.SearchByBrandNameAndItem(pro.Brand_Name, pro.Item)
	if err != nil {
		return nil, nil, errors.New("error with server")
	}
	if name.Brand_Name != "" && name.Item != "" {
		return nil, nil, errors.New("entered product with this same brand already exists")
	}
	category, err := pu.productRepo.GetCategoryByID(pro.CategoryID)
	if err != nil {
		return nil, nil, errors.New("error with server")
	}
	NewProduct := &entity.Product{
		Brand_Name:  pro.Brand_Name,
		Description: pro.Description,
		Price:       pro.Price,
		Item:        pro.Item,
		CategoryID:  pro.CategoryID,
		Quantity:    pro.Quantity,
		ImageURL:    pro.ImageURL,
	}
	CreatedProduct, err := pu.productRepo.CreateProduct(NewProduct)
	if err != nil {
		return nil, nil, err
	}
	return CreatedProduct, category, nil
}

func (pu *ProductUsecase) ExecutePaginatedProducts(offset int, limit int) ([]entity.Product, error) {
	products, err := pu.productRepo.GetPaginatedProducts(offset, limit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (pu *ProductUsecase) GetProductByID(ID uint) (*entity.Product, *entity.Category, error) {
	product, err := pu.productRepo.GetProductByID(ID)
	if err != nil {
		return nil, nil, err
	}
	category, err := pu.productRepo.GetCategoryByID(product.CategoryID)
	if err != nil {
		return nil, nil, errors.New("error with server")
	}
	return product, category, nil
}

func (pu *ProductUsecase) ProductUpdate(product *utils.EditProductRequest) (*entity.Product, *entity.Category, error) {
	var UpdatedProduct *entity.Product
	ExistingProduct, err := pu.productRepo.GetProductByID(product.ID)
	if err != nil {
		return nil, nil, errors.New("failed to fetch  product")
	}
	if product.Brand_Name != nil {
		ExistingProduct.Brand_Name = *product.Brand_Name
	}
	if product.Description != nil {
		ExistingProduct.Description = *product.Description
	}
	if product.Price != nil {
		ExistingProduct.Price = *product.Price
	}
	if product.Quantity != nil {
		ExistingProduct.Quantity = *product.Quantity
	}
	if product.Category != nil {
		ExistingProduct.CategoryID = *product.Category
	}
	UpdatedProduct, err1 := pu.productRepo.ProductUpdate(ExistingProduct)
	if err1 != nil {
		return nil, nil, err1
	}
	category, err := pu.productRepo.GetCategoryByID(UpdatedProduct.CategoryID)
	if err != nil {
		return nil, nil, errors.New("can't fetch category of the product")
	}
	return UpdatedProduct, category, nil
}

func (pu *ProductUsecase) DeleteProduct(id uint) error {
	err := pu.productRepo.ProductDelete(id)
	if err != nil {
		return errors.New("can't execute deletion")
	}
	return nil
}

func (pu *ProductUsecase) FetchProductByBrandName(BrandName string, offset int, limit int) ([]entity.Product, error) {
	FilteredProducts, err := pu.productRepo.GetPaginatedProductsByBrandName(BrandName, offset, limit)
	if err != nil {
		return nil, errors.New("can't execute fetching products")
	}
	return FilteredProducts, nil
}

func (pu *ProductUsecase) FetchBrandName(BrandName string) (*entity.Product, error) {
	FetchedBrand, err := pu.productRepo.GetByBrand(BrandName)
	if err != nil {
		return nil, errors.New("can't fetch brand")
	}
	return FetchedBrand, nil
}
func (pu *ProductUsecase) FetchItemName(ItemName string) (*entity.Product, error) {
	FetchedItem, err := pu.productRepo.GetByItem(ItemName)
	if err != nil {
		return nil, errors.New("can't fetch brand")
	}
	return FetchedItem, nil
}

func (pu *ProductUsecase) FetchPaginatedProductByItemName(ItemName string, offset int, limit int) ([]entity.Product, error) {
	FetchedProducts, err := pu.productRepo.GetPaginatedProductsByItemName(ItemName, offset, limit)
	if err != nil {
		return nil, errors.New("can't fetch products")
	}
	return FetchedProducts, nil
}

func (pu *ProductUsecase) ExecuteTotalProducts() (int, error) {
	TotalOfProducts, err := pu.productRepo.GetTotalOfProducts()
	if err != nil {
		return 0, errors.New("can't get the total of products")
	}
	return TotalOfProducts, nil
}

func (pu *ProductUsecase) ExecuteTotalOfProductsByBrand(BrandName string) (int, error) {
	TotalOfProductsByBrand, err := pu.productRepo.GetTotalOfProductsByBrand(BrandName)
	if err != nil {
		return 0, errors.New("can fetch the total of products in this brand")
	}
	return TotalOfProductsByBrand, nil
}

func (pu *ProductUsecase) ExecuteTotalOfProductsByItemName(ItemName string) (int, error) {
	TotalOfProductsByItemName, err := pu.productRepo.GetTotalOfProductsByItemName(ItemName)
	if err != nil {
		return 0, errors.New("can't fetch the total of products in this item name")
	}
	return TotalOfProductsByItemName, nil
}
