// Package database opens the application's Postgres connection pool.
package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const connectTimeout = 5 * time.Second

// Connect opens a connection pool and checks that Postgres is reachable, so a
// broken setup fails at startup rather than on the first request.
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set: copy .env.example to .env, or set it in the environment")
	}

	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("configuring database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connecting to database (is Postgres running? try make db-up): %w", err)
	}

	slog.Info("connected to database")
	return pool, nil
}
