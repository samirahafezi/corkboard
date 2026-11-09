# Migrations

Numbered SQL files that build the database schema. They run in order, and the `goose_db_version` table records which ones have been applied.

- Create one with `make migration name=create_users`, then fill in its `-- +goose Up` and `-- +goose Down` sections.
- Apply pending migrations with `make migrate`. Check them with `make migrate-status`, and undo the latest with `make migrate-down`.
- Never edit a migration after it has been merged. Add a new one instead.
- Each migration runs in a transaction, so a failure leaves nothing half-applied.
