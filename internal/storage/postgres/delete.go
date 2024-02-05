package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"webinar-testing/pkg/models"
)

func (s *storage) Delete(ctx context.Context, order models.Order) error {
	newCtx, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	b := &pgx.Batch{}

	for good, qty := range order.Goods {
		b.Queue("INSERT INTO orders (user_id, good, quantity) VALUES ($1, $2, $3) ", order.UserID, good, -qty)
	}

	if err := s.pool.SendBatch(newCtx, b).Close(); err != nil {
		return fmt.Errorf("postgres delete batch sent: %w", err)
	}

	return nil
}

func (s *storage) DeleteAllByUser(ctx context.Context, id models.UserID) error {
	newCtx, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := s.pool.Exec(newCtx, "DELETE FROM orders WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("postgres delete: %w", err)
	}

	return nil
}
