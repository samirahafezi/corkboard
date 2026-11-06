// Command server runs the Corkboard HTTP API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/samirahafezi/corkboard/api/internal/api"
)

const (
	defaultPort = "8080"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second

	shutdownTimeout = 10 * time.Second

	dbConnectTimeout = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

// run starts the server and blocks until it fails or receives a stop signal.
func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := connectDB(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	// Deferred calls run in reverse order, so the pool closes after the
	// server has finished shutting down.
	defer pool.Close()

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           api.Handler(api.NewStrictHandler(api.Server{}, nil)),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	// The buffer lets the goroutine send and exit even if nothing is receiving.
	serveErr := make(chan error, 1)
	go func() {
		slog.Info("server listening", "port", port)
		serveErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		return fmt.Errorf("serving: %w", err)
	case <-ctx.Done():
		slog.Info("shutting down")
	}

	// The signal context is already canceled, so shutdown needs a fresh one.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutting down: %w", err)
	}

	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serving: %w", err)
	}

	slog.Info("server stopped cleanly")
	return nil
}

// connectDB opens a connection pool and checks that Postgres is reachable,
// so a broken setup fails at startup rather than on the first request.
func connectDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set: copy .env.example to .env, or set it in the environment")
	}

	ctx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
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
