// Command migrate applies, undoes, and creates database migrations.
//
// Usage:
//
//	migrate up            apply all pending migrations
//	migrate down          undo the most recent migration
//	migrate status        list migrations and whether they've been applied
//	migrate create NAME   create a new, empty SQL migration
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Registers the "pgx" driver with database/sql.
	"github.com/pressly/goose/v3"
)

// migrationsDir is relative to api/, where the Makefile runs this command.
const migrationsDir = "migrations"

const usage = "usage: migrate up | down | status | create NAME"

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		slog.Error("migrate failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New(usage)
	}

	switch args[0] {
	case "create":
		// Creating a file doesn't need a database connection.
		if len(args) != 2 {
			return errors.New("usage: migrate create NAME")
		}
		goose.SetSequential(true)
		return goose.Create(nil, migrationsDir, args[1], "sql")
	case "up", "down", "status":
	default:
		return fmt.Errorf("unknown command %q: %s", args[0], usage)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is not set: copy .env.example to .env, or set it in the environment")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("configuring database: %w", err)
	}
	defer func() { _ = db.Close() }()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS(migrationsDir))
	if errors.Is(err, goose.ErrNoMigrations) {
		slog.Info("no migrations yet", "dir", migrationsDir)
		return nil
	}
	if err != nil {
		return fmt.Errorf("loading migrations: %w", err)
	}

	switch args[0] {
	case "up":
		results, err := provider.Up(ctx)
		for _, r := range results {
			slog.Info("applied", "migration", r.Source.Path, "duration", r.Duration)
		}
		if err != nil {
			return fmt.Errorf("applying migrations: %w", err)
		}
		if len(results) == 0 {
			slog.Info("already up to date")
		}
	case "down":
		result, err := provider.Down(ctx)
		if err != nil {
			return fmt.Errorf("undoing migration: %w", err)
		}
		slog.Info("undone", "migration", result.Source.Path)
	case "status":
		statuses, err := provider.Status(ctx)
		if err != nil {
			return fmt.Errorf("reading status: %w", err)
		}
		for _, s := range statuses {
			slog.Info("migration", "file", s.Source.Path, "state", s.State)
		}
	}
	return nil
}
