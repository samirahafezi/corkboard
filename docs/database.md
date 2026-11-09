# Local database

How to run, inspect, and change the local Postgres database.

Run every command from the repo root. `docker compose` and `make` both look for their files in the current folder.

## Start and stop

| Task | Command |
| --- | --- |
| Start Postgres, waiting until it's ready | `make db-up` |
| Stop Postgres and keep the data | `make db-down` |
| Check whether it's running | `docker compose ps` |
| See its logs | `docker compose logs postgres` |
| Follow its logs live (Ctrl+C to stop) | `docker compose logs -f postgres` |

Postgres has to be running before `make dev-api` or any `make migrate...` command.

## Open a SQL prompt (psql)

```sh
docker compose exec postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
```

This opens a `corkboard=#` prompt. SQL statements end with `;`, and backslash commands don't.

| Command | What it does |
| --- | --- |
| `\dt` | List tables |
| `\d users` | Show a table's columns, indexes, and rules |
| `\x` | Toggle expanded display, which is easier for wide rows |
| `\?` | List all psql commands |
| `\q` | Quit |

To run a single query without opening the prompt:

```sh
docker compose exec postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "SELECT * FROM users;"'
```

## Handy queries

```sql
-- All users, newest first
SELECT * FROM users ORDER BY created_at DESC;

-- Find a user by email, ignoring case
SELECT * FROM users WHERE lower(email) = lower('Sam@School.org');

-- Add a user by hand
INSERT INTO users (email, role) VALUES ('you@school.org', 'admin') RETURNING *;

-- Remove all users
DELETE FROM users;

-- Which migrations have run
SELECT version_id, is_applied, tstamp FROM goose_db_version ORDER BY id;
```

The allowed roles are `student`, `staff`, and `admin`.

## Migrations

Migrations live in `api/migrations/`. See the README there for conventions.

| Task | Command |
| --- | --- |
| Apply pending migrations | `make migrate` |
| See which have run | `make migrate-status` |
| Undo the most recent one | `make migrate-down` |
| Create a new one | `make migration name=create_clubs` |

Never edit a migration after it's been merged. Add a new one instead.

## Queries (sqlc)

- SQL queries live in `api/queries/`, and sqlc generates Go for them in `api/internal/db/`.
- After changing a query or a migration, run `make generate`, then commit the generated files.
- `make check-generated` fails if you forgot. CI runs it too.

## Start over with an empty database

This deletes all local data:

```sh
docker compose down -v   # stop Postgres and delete its data volume
make db-up               # start a fresh, empty Postgres
make migrate             # recreate the tables
```

## Connect with a GUI app

Any Postgres client works, such as TablePlus, DBeaver, Postico, or pgAdmin.

| Setting | Value |
| --- | --- |
| Host | `127.0.0.1` (not `localhost`) |
| Port | `5432` |
| Database | `corkboard` |
| User | `corkboard` |
| Password | `POSTGRES_PASSWORD` from `.env` |
| SSL | Off / disable |

Postgres only accepts connections from your own machine.

## Troubleshooting

| Error | Fix |
| --- | --- |
| `no configuration file provided: not found` | You're not in the repo root. `cd` there first. |
| `service "postgres" is not running` | Run `make db-up` |
| `required variable POSTGRES_USER is missing a value` | `.env` is missing. Run `cp .env.example .env` |
| `make: *** No rule to make target` | You're in a subfolder like `api/`. Go back to the repo root. |
| API says `is Postgres running? try make db-up` | Postgres is stopped, or `DATABASE_URL` in `.env` doesn't match the Postgres settings |
| `docker: command not found` | Docker Desktop isn't running, or its `bin` folder isn't on your `PATH` |
| Port 5432 already in use | Another Postgres is running on your machine. Stop it, or change the port in `compose.yaml` and `DATABASE_URL`. |
