package models

import (
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

func (r *CategoriesRepositoryImpl) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepositoryImpl) CreateCategory(category *Category) error {
	return r.db.Create(category).Error
}
