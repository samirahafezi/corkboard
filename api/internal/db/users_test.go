package db_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/samirahafezi/corkboard/api/internal/db"
	"github.com/samirahafezi/corkboard/api/internal/dbtest"
)

// Postgres error codes. See https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	uniqueViolation = "23505"
	checkViolation  = "23514"
)

func TestMain(m *testing.M) { dbtest.Main(m) }

func TestCreateAndGetUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	q := db.New(dbtest.NewPool(t))

	created, err := q.CreateUser(ctx, db.CreateUserParams{Email: "Sam@School.org", Role: "student"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if created.ID == 0 {
		t.Error("ID was not set")
	}
	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
	if created.Email != "Sam@School.org" || created.Role != "student" {
		t.Errorf("created = %+v, want email Sam@School.org and role student", created)
	}

	byEmail, err := q.GetUserByEmail(ctx, "sam@SCHOOL.org")
	if err != nil {
		t.Fatalf("GetUserByEmail with different case: %v", err)
	}
	if byEmail.ID != created.ID {
		t.Errorf("GetUserByEmail ID = %d, want %d", byEmail.ID, created.ID)
	}

	byID, err := q.GetUserByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if byID != created {
		t.Errorf("GetUserByID = %+v, want %+v", byID, created)
	}
}

func TestGetUserNotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	q := db.New(dbtest.NewPool(t))

	if _, err := q.GetUserByEmail(ctx, "nobody@school.org"); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("GetUserByEmail error = %v, want pgx.ErrNoRows", err)
	}
	if _, err := q.GetUserByID(ctx, 12345); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("GetUserByID error = %v, want pgx.ErrNoRows", err)
	}
}

func TestCreateUserRejectsInvalidUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		email    string
		role     string
		wantCode string
	}{
		{name: "duplicate email in a different case", email: "EXISTING@school.org", role: "student", wantCode: uniqueViolation},
		{name: "unknown role", email: "lead@school.org", role: "club_lead", wantCode: checkViolation},
		{name: "empty email", email: "", role: "student", wantCode: checkViolation},
		{name: "email over 254 characters", email: strings.Repeat("a", 244) + "@school.org", role: "student", wantCode: checkViolation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			q := db.New(dbtest.NewPool(t))

			if _, err := q.CreateUser(ctx, db.CreateUserParams{Email: "existing@school.org", Role: "student"}); err != nil {
				t.Fatalf("creating existing user: %v", err)
			}

			_, err := q.CreateUser(ctx, db.CreateUserParams{Email: tt.email, Role: tt.role})

			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				t.Fatalf("error = %v, want a Postgres error with code %s", err, tt.wantCode)
			}
			if pgErr.Code != tt.wantCode {
				t.Errorf("error code = %s (%s), want %s", pgErr.Code, pgErr.Message, tt.wantCode)
			}
		})
	}
}
