package category

import (
	"context"

	"coi/internal/model"
)

type CategoryService interface {
	GetCategoryByID(ctx context.Context, id int) (*model.Category, error)
	ListCategories(ctx context.Context) ([]model.Category, error)
}

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *categoryService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.ListActive(ctx)
}
