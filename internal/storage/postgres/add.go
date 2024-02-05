package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"

	"webinar-testing/pkg/errs"
	"webinar-testing/pkg/models"
)

func (s *storage) Add(ctx context.Context, order models.Order) error {
	newCtx, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	b := &pgx.Batch{}

	for good, qty := range order.Goods {
		b.Queue("INSERT INTO orders (user_id, good, quantity) VALUES ($1, $2, $3)", order.UserID, good, qty)
	}

	if err := s.pool.SendBatch(newCtx, b).Close(); err != nil {
		log.Println("postgres connect:", err)
		return errs.NewWithError("send batch", err)
	}

	return nil
}
