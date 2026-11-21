// Command create-admin adds a user with the admin role, so the first admin
// can sign in and add everyone else.
//
// Usage:
//
//	create-admin -email you@school.org
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/samirahafezi/corkboard/api/internal/database"
	"github.com/samirahafezi/corkboard/api/internal/db"
	"github.com/samirahafezi/corkboard/api/internal/users"
)

func main() {
	email := flag.String("email", "", "email address of the new admin")
	flag.Parse()

	if *email == "" {
		fmt.Fprintln(os.Stderr, "usage: create-admin -email you@school.org")
		os.Exit(2)
	}

	if err := run(context.Background(), *email); err != nil {
		slog.Error("create-admin failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, email string) error {
	pool, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	admin, err := users.CreateAdmin(ctx, db.New(pool), email)
	if errors.Is(err, users.ErrUserExists) {
		return fmt.Errorf("%w. Their role wasn't changed", err)
	}
	if err != nil {
		return err
	}

	slog.Info("created admin", "id", admin.ID, "email", admin.Email)
	return nil
}
