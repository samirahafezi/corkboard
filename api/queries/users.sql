-- name: CreateUser :one
INSERT INTO users (email, role)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByEmail :one
-- Matches the case-insensitive unique index on lower(email).
SELECT * FROM users
WHERE lower(email) = lower(sqlc.arg(email));

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;
