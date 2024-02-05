package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"webinar-testing/pkg/models"
)

func (s *storage) ListByUser(ctx context.Context, id models.UserID) (order models.Order, err error) {
	newCtx, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	rows, err := s.pool.Query(newCtx, "SELECT (id, user_id, good, quantity, created_at) FROM orders WHERE user_id = $1", id)
	if err != nil {
		return order, fmt.Errorf("postgres list query: %w", err)
	}

	res, err := pgx.CollectRows(rows, pgx.RowTo[ListResult])
	if err != nil {
		return order, fmt.Errorf("postgres list collect rows: %w", err)
	}

	return listResultsToOrder(res, id), nil
}
