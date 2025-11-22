package services

import (
	"context"
	"errors"

	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/repositories"
)

type CategoryService interface {
	Create(ctx context.Context, userID, name string, parentID *string) (*models.Category, error)
	List(ctx context.Context, userID string, parentID *string) ([]models.Category, error)
	Get(ctx context.Context, userID, id string) (*models.Category, error)
	Update(ctx context.Context, userID, id, name string, parentID *string) (*models.Category, error)
	Delete(ctx context.Context, userID, id string) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, userID, name string, parentID *string) (*models.Category, error) {
	if parentID != nil {
		parent, err := s.repo.FindByID(ctx, *parentID, userID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, errors.New("parent category not found")
		}
	}

	cat := &models.Category{
		UserID:   userID,
		Name:     name,
		ParentID: parentID,
	}
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) List(ctx context.Context, userID string, parentID *string) ([]models.Category, error) {
	return s.repo.List(ctx, userID, parentID)
}

func (s *categoryService) Get(ctx context.Context, userID, id string) (*models.Category, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *categoryService) Update(ctx context.Context, userID, id, name string, parentID *string) (*models.Category, error) {
	cat, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("category not found")
	}

	if parentID != nil {
		parent, err := s.repo.FindByID(ctx, *parentID, userID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, errors.New("parent category not found")
		}
	}

	cat.Name = name
	cat.ParentID = parentID

	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Delete(ctx context.Context, userID, id string) error {
	// NOTE: here you could check for transactions referencing this category and either block or nullify.
	return s.repo.Delete(ctx, id, userID)
}
