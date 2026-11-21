package users_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/samirahafezi/corkboard/api/internal/db"
	"github.com/samirahafezi/corkboard/api/internal/dbtest"
	"github.com/samirahafezi/corkboard/api/internal/users"
)

func TestMain(m *testing.M) { dbtest.Main(m) }

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()

	valid := []struct {
		in   string
		want string
	}{
		{in: "sam@school.org", want: "sam@school.org"},
		{in: "Sam@School.ORG", want: "sam@school.org"},
		{in: "  sam@school.org\n", want: "sam@school.org"},
		{in: "sam.lee+clubs@mail.school.org", want: "sam.lee+clubs@mail.school.org"},
	}
	for _, tt := range valid {
		got, err := users.NormalizeEmail(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("NormalizeEmail(%q) = %q, %v; want %q, nil", tt.in, got, err, tt.want)
		}
	}

	invalid := []string{
		"",
		"   ",
		"sam",
		"sam@",
		"@school.org",
		"sam@school",
		"sam@school.",
		"Sam <sam@school.org>",
		"<sam@school.org>",
		"sam@school.org, lee@school.org",
		"sam @school.org",
		strings.Repeat("a", 244) + "@school.org",
	}
	for _, in := range invalid {
		if got, err := users.NormalizeEmail(in); !errors.Is(err, users.ErrInvalidEmail) {
			t.Errorf("NormalizeEmail(%q) = %q, %v; want ErrInvalidEmail", in, got, err)
		}
	}
}

func TestCreateAdmin(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	q := db.New(dbtest.NewPool(t))

	admin, err := users.CreateAdmin(ctx, q, " Admin@School.org ")
	if err != nil {
		t.Fatalf("CreateAdmin: %v", err)
	}
	if admin.Email != "admin@school.org" || admin.Role != users.RoleAdmin {
		t.Errorf("created %+v, want email admin@school.org and role admin", admin)
	}

	if _, err := users.CreateAdmin(ctx, q, "ADMIN@school.org"); !errors.Is(err, users.ErrUserExists) {
		t.Errorf("second CreateAdmin error = %v, want ErrUserExists", err)
	}

	if _, err := users.CreateAdmin(ctx, q, "not-an-email"); !errors.Is(err, users.ErrInvalidEmail) {
		t.Errorf("CreateAdmin with bad email error = %v, want ErrInvalidEmail", err)
	}
}

func TestCreateAdminRefusesExistingStudent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	q := db.New(dbtest.NewPool(t))

	if _, err := q.CreateUser(ctx, db.CreateUserParams{Email: "sam@school.org", Role: users.RoleStudent}); err != nil {
		t.Fatalf("creating student: %v", err)
	}

	if _, err := users.CreateAdmin(ctx, q, "sam@school.org"); !errors.Is(err, users.ErrUserExists) {
		t.Fatalf("CreateAdmin error = %v, want ErrUserExists", err)
	}

	student, err := q.GetUserByEmail(ctx, "sam@school.org")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	if student.Role != users.RoleStudent {
		t.Errorf("role = %q, want it to stay %q", student.Role, users.RoleStudent)
	}
}
