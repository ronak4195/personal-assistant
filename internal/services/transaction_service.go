package services

import (
	"context"
	"errors"
	"time"

	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/repositories"
)

type TransactionService interface {
	Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error)
	Get(ctx context.Context, userID, id string) (*models.Transaction, error)
	List(ctx context.Context, f repositories.TransactionFilter) ([]models.Transaction, int64, error)
	Update(ctx context.Context, userID string, tx *models.Transaction) (*models.Transaction, error)
	Delete(ctx context.Context, userID, id string) error
}

type transactionService struct {
	repo         repositories.TransactionRepository
	categoryRepo repositories.CategoryRepository
}

func NewTransactionService(repo repositories.TransactionRepository, catRepo repositories.CategoryRepository) TransactionService {
	return &transactionService{
		repo:         repo,
		categoryRepo: catRepo,
	}
}

func (s *transactionService) Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error) {
	if tx.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if tx.Type != models.TransactionTypeIncome && tx.Type != models.TransactionTypeExpense {
		return nil, errors.New("invalid type")
	}

	if tx.CategoryID != nil {
		if _, err := s.categoryRepo.FindByID(ctx, *tx.CategoryID, tx.UserID); err != nil {
			return nil, err
		}
	}
	if tx.SubcategoryID != nil {
		if _, err := s.categoryRepo.FindByID(ctx, *tx.SubcategoryID, tx.UserID); err != nil {
			return nil, err
		}
	}
	if tx.Date.IsZero() {
		tx.Date = time.Now().UTC()
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *transactionService) Get(ctx context.Context, userID, id string) (*models.Transaction, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *transactionService) List(ctx context.Context, f repositories.TransactionFilter) ([]models.Transaction, int64, error) {
	return s.repo.List(ctx, f)
}

func (s *transactionService) Update(ctx context.Context, userID string, tx *models.Transaction) (*models.Transaction, error) {
	existing, err := s.repo.FindByID(ctx, tx.ID, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("transaction not found")
	}

	// Update allowed fields
	if tx.Type != "" {
		if tx.Type != models.TransactionTypeIncome && tx.Type != models.TransactionTypeExpense {
			return nil, errors.New("invalid type")
		}
		existing.Type = tx.Type
	}
	if tx.Amount > 0 {
		existing.Amount = tx.Amount
	}
	if tx.Currency != "" {
		existing.Currency = tx.Currency
	}
	if tx.CategoryID != nil {
		existing.CategoryID = tx.CategoryID
	}
	if tx.SubcategoryID != nil {
		existing.SubcategoryID = tx.SubcategoryID
	}
	if tx.Description != nil {
		existing.Description = tx.Description
	}
	if !tx.Date.IsZero() {
		existing.Date = tx.Date
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *transactionService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}
