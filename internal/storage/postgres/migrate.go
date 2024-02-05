package postgres

import (
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations
var migrations embed.FS

func migrate(pool *pgxpool.Pool, version int64) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate set dialect postgres: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.UpTo(db, "migrations", version); err != nil {
		return fmt.Errorf("postgres migrate up: %w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("postgres migrate close db: %w", err)
	}
	return nil
}
