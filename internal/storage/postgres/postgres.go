package postgres

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5/pgxpool"

	"webinar-testing/pkg/models"
)

type Storager interface {
	Add(ctx context.Context, goods models.Order) error
	ListByUser(ctx context.Context, id models.UserID) (models.Order, error)
	Delete(ctx context.Context, goods models.Order) error
	DeleteAllByUser(ctx context.Context, id models.UserID) error
	io.Closer
}

var _ Storager = &storage{}

type storage struct {
	cfg  *Config
	pool *pgxpool.Pool
	do   func(ctx context.Context)
}

func (s *storage) Close() error {
	s.pool.Close()
	return nil
}

func New(cfg *Config) (*storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.connectionString())
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	if err := migrate(pool, cfg.MigrationVersion); err != nil {
		return nil, err
	}

	return &storage{cfg: cfg, pool: pool}, nil
}
