// Package dbtest gives tests a real, migrated Postgres database.
//
// Each test package starts one Postgres container in TestMain:
//
//	func TestMain(m *testing.M) { dbtest.Main(m) }
//
// Each test then calls NewPool to get its own empty, fully migrated database,
// which is deleted when the test finishes. Tests can't see each other's data,
// so they can safely run in parallel.
package dbtest

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Registers the "pgx" driver for goose.
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/samirahafezi/corkboard/api/migrations"
)

// image matches the version in compose.yaml, so tests run against the same Postgres.
const image = "postgres:18-alpine"

// templateDB is migrated once. Each test's database is a copy of it.
const templateDB = "corkboard_template"

const setupTimeout = 2 * time.Minute

var (
	// admin is connected to the "postgres" database and creates and drops test databases.
	admin *pgxpool.Pool
	// baseConfig points at the container; tests swap in their own database name.
	baseConfig *pgxpool.Config
	// dbCount gives each test database a unique name.
	dbCount atomic.Int64
	// cloneMu serializes copies, because Postgres can't copy a template
	// while another copy of it is in progress.
	cloneMu sync.Mutex
)

// Main starts Postgres, applies the migrations, runs the package's tests, and
// stops the container. Call it from TestMain.
func Main(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), setupTimeout)
	container, err := setup(ctx)
	cancel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "dbtest: %v\nThese tests need Docker. Is Docker Desktop running?\n", err)
		os.Exit(1)
	}

	code := m.Run()

	admin.Close()
	if err := testcontainers.TerminateContainer(container); err != nil {
		fmt.Fprintf(os.Stderr, "dbtest: stopping container: %v\n", err)
	}
	os.Exit(code)
}

func setup(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(ctx, image,
		postgres.WithDatabase(templateDB),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("starting Postgres container: %w", err)
	}

	templateURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return container, fmt.Errorf("reading connection string: %w", err)
	}

	if err := migrate(ctx, templateURL); err != nil {
		return container, err
	}

	baseConfig, err = pgxpool.ParseConfig(templateURL)
	if err != nil {
		return container, fmt.Errorf("parsing connection string: %w", err)
	}

	adminConfig := baseConfig.Copy()
	adminConfig.ConnConfig.Database = "postgres"
	admin, err = pgxpool.NewWithConfig(ctx, adminConfig)
	if err != nil {
		return container, fmt.Errorf("connecting as admin: %w", err)
	}
	return container, nil
}

// migrate applies every migration to the template database, then disconnects,
// because Postgres can only copy a database that has no open connections.
func migrate(ctx context.Context, url string) error {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return fmt.Errorf("opening template database: %w", err)
	}
	defer func() { _ = db.Close() }()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return fmt.Errorf("loading migrations: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}
	return nil
}

// NewPool returns a connection pool for a new, empty, fully migrated database.
// The database is deleted when the test finishes.
func NewPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if admin == nil {
		t.Fatal("dbtest: call dbtest.Main(m) from TestMain first")
	}
	ctx := context.Background()

	name := fmt.Sprintf("test_%d", dbCount.Add(1))
	quoted := pgx.Identifier{name}.Sanitize()

	cloneMu.Lock()
	_, err := admin.Exec(ctx, "CREATE DATABASE "+quoted+" TEMPLATE "+pgx.Identifier{templateDB}.Sanitize())
	cloneMu.Unlock()
	if err != nil {
		t.Fatalf("dbtest: creating test database: %v", err)
	}

	config := baseConfig.Copy()
	config.ConnConfig.Database = name
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("dbtest: connecting to test database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		if _, err := admin.Exec(context.Background(), "DROP DATABASE "+quoted+" WITH (FORCE)"); err != nil {
			t.Errorf("dbtest: dropping %s: %v", strings.Trim(quoted, `"`), err)
		}
	})
	return pool
}
