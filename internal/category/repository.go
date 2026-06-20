package category

import (
	"context"

	"coi/internal/model"
	"coi/pkg/txmanager"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindByID(ctx context.Context, id int) (*model.Category, error)
	ListActive(ctx context.Context) ([]model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindByID(ctx context.Context, id int) (*model.Category, error) {
	db := txmanager.GetTx(ctx, r.db)
	var category model.Category
	err := db.WithContext(ctx).Where("is_active = ?", true).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) ListActive(ctx context.Context) ([]model.Category, error) {
	db := txmanager.GetTx(ctx, r.db)
	var categories []model.Category
	err := db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&categories).Error
	return categories, err
}
