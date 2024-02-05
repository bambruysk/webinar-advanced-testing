package service

import (
	"context"
	"fmt"
	"log"

	"webinar-testing/pkg/models"
)

type Storager interface {
	Add(ctx context.Context, goods models.Order) error
	ListByUser(ctx context.Context, id models.UserID) (models.Order, error)
	Delete(ctx context.Context, goods models.Order) error
	DeleteAllByUser(ctx context.Context, id models.UserID) error
}

type service struct {
	storage Storager
}

func New(storage Storager) *service {
	return &service{storage: storage}
}

func (s *service) Add(ctx context.Context, goods models.Order) error {
	if err := goods.Validate(); err != nil {
		return fmt.Errorf("service add: %w", err)
	}

	if err := s.storage.Add(ctx, goods); err != nil {
		log.Println("postgres connect:", err)
		return fmt.Errorf("service add: %w", err)
	}

	return nil
}

func (s *service) ListByUser(ctx context.Context, id models.UserID) (order models.Order, err error) {
	order, err = s.storage.ListByUser(ctx, id)
	if err != nil {
		return order, fmt.Errorf("service list by user: %w", err)
	}

	return order, nil
}

func (s *service) Delete(ctx context.Context, goods models.Order) (err error) {
	if err = s.storage.Delete(ctx, goods); err != nil {
		return fmt.Errorf("service delete: %w", err)
	}

	return nil
}

func (s *service) DeleteAllByUser(ctx context.Context, id models.UserID) error {
	if err := s.storage.DeleteAllByUser(ctx, id); err != nil {
		return fmt.Errorf("service delete all by user: %w", err)
	}

	return nil
}
