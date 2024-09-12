package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPgx(ctx context.Context, connectionString string) (*sqlx.DB, error) {
	pool, err := NewPgxPool(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create to pgx pool: %w", err)
	}

	return WrapPgxPool(ctx, pool)
}

func NewPgxPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	if poolConfig == nil {
		return nil, errors.New("parsed config is nil")
	}

	return pgxpool.NewWithConfig(ctx, poolConfig)
}

func WrapPgxPool(ctx context.Context, pool *pgxpool.Pool) (*sqlx.DB, error) {
	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), "pgx")
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
