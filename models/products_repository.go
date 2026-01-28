package models

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductsRepositoryImpl struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepositoryImpl {
	return &ProductsRepositoryImpl{
		db: db,
	}
}

func (r *ProductsRepositoryImpl) GetAllProducts(ctx context.Context, offset, limit int, categoryCode string, priceLessThan *decimal.Decimal) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.WithContext(ctx).Model(&Product{}).Preload("Category").Preload("Variants")

	if categoryCode != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", categoryCode)
	}

	if priceLessThan != nil {
		query = query.Where("products.price < ?", priceLessThan)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepositoryImpl) GetProductByCode(ctx context.Context, code string) (*Product, error) {
	var product Product
	if err := r.db.WithContext(ctx).Preload("Category").Preload("Variants").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
