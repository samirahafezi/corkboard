// Package users holds the rules for user accounts: valid email addresses,
// roles, and creating users.
package users

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/samirahafezi/corkboard/api/internal/db"
)

// Base roles, matching the CHECK constraint on users.role.
const (
	RoleStudent = "student"
	RoleStaff   = "staff"
	RoleAdmin   = "admin"
)

// maxEmailLength matches the CHECK constraint on users.email.
const maxEmailLength = 254

// uniqueViolation is the Postgres error code for a duplicate key.
const uniqueViolation = "23505"

var (
	// ErrInvalidEmail means the input isn't a single, plain email address.
	ErrInvalidEmail = errors.New("not a valid email address")
	// ErrUserExists means a user with that email already exists.
	ErrUserExists = errors.New("a user with that email already exists")
)

// NormalizeEmail checks that s is a single, plain email address, such as
// sam@school.org, and returns it trimmed and lowercased.
func NormalizeEmail(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" || len(s) > maxEmailLength {
		return "", ErrInvalidEmail
	}

	addr, err := mail.ParseAddress(s)
	// ParseAddress also accepts forms like "Sam <sam@school.org>". Requiring
	// the parsed address to equal the input rejects names and angle brackets.
	if err != nil || addr.Name != "" || addr.Address != s {
		return "", ErrInvalidEmail
	}

	// Require a dot in the domain, so typos like sam@school are caught.
	_, domain, _ := strings.Cut(addr.Address, "@")
	if !strings.Contains(domain, ".") || strings.HasSuffix(domain, ".") {
		return "", ErrInvalidEmail
	}

	return strings.ToLower(addr.Address), nil
}

// CreateAdmin creates a user with the admin role. It returns ErrInvalidEmail
// or ErrUserExists for those cases.
func CreateAdmin(ctx context.Context, q *db.Queries, email string) (db.User, error) {
	normalized, err := NormalizeEmail(email)
	if err != nil {
		return db.User{}, err
	}

	user, err := q.CreateUser(ctx, db.CreateUserParams{Email: normalized, Role: RoleAdmin})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
		return db.User{}, ErrUserExists
	}
	if err != nil {
		return db.User{}, fmt.Errorf("creating admin: %w", err)
	}
	return user, nil
}
