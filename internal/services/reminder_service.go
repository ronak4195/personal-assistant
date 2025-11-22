package services

import (
	"context"
	"errors"
	"time"

	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/repositories"
)

type ReminderService interface {
	Create(ctx context.Context, r *models.Reminder) (*models.Reminder, error)
	List(ctx context.Context, f repositories.ReminderFilter) ([]models.Reminder, error)
	Get(ctx context.Context, userID, id string) (*models.Reminder, error)
	Update(ctx context.Context, userID string, r *models.Reminder) (*models.Reminder, error)
	Delete(ctx context.Context, userID, id string) error
}

type reminderService struct {
	repo repositories.ReminderRepository
}

func NewReminderService(repo repositories.ReminderRepository) ReminderService {
	return &reminderService{repo: repo}
}

func (s *reminderService) Create(ctx context.Context, r *models.Reminder) (*models.Reminder, error) {
	if r.Title == "" {
		return nil, errors.New("title is required")
	}
	if r.DueAt.IsZero() {
		r.DueAt = time.Now().Add(1 * time.Hour)
	}
	if r.RepeatInterval == "" {
		r.RepeatInterval = models.RepeatNone
	}
	r.IsActive = true

	if err := s.repo.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *reminderService) List(ctx context.Context, f repositories.ReminderFilter) ([]models.Reminder, error) {
	return s.repo.List(ctx, f)
}

func (s *reminderService) Get(ctx context.Context, userID, id string) (*models.Reminder, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *reminderService) Update(ctx context.Context, userID string, r *models.Reminder) (*models.Reminder, error) {
	existing, err := s.repo.FindByID(ctx, r.ID, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("reminder not found")
	}

	if r.Title != "" {
		existing.Title = r.Title
	}
	if r.Description != nil {
		existing.Description = r.Description
	}
	if !r.DueAt.IsZero() {
		existing.DueAt = r.DueAt
	}
	if r.RepeatInterval != "" {
		existing.RepeatInterval = r.RepeatInterval
	}
	existing.IsActive = r.IsActive

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *reminderService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}
