package models

import (
	"context"

	"gorm.io/gorm"
)

type CategoriesRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepositoryImpl {
	return &CategoriesRepositoryImpl{
		db: db,
	}
}

func (r *CategoriesRepositoryImpl) GetAllCategories(ctx context.Context) ([]Category, error) {
	var categories []Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepositoryImpl) CreateCategory(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}
