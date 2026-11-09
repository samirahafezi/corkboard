-- +goose Up
-- Users are created by an admin. Only the minimum needed to sign in and
-- authorize is stored: an email address and a base role.
CREATE TABLE users (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email      text NOT NULL CHECK (email <> '' AND length(email) <= 254),
    role       text NOT NULL CHECK (role IN ('student', 'staff', 'admin')),
    created_at timestamptz NOT NULL DEFAULT now()
);

-- Emails are unique regardless of letter case.
CREATE UNIQUE INDEX users_email_key ON users (lower(email));

-- +goose Down
DROP TABLE users;
