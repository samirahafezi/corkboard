# Decisions

A short log of meaningful technical choices, newest first. The initial stack is listed in `CLAUDE.md`.

## 2026-11-08: Run goose migrations through our own small command

**Decision:** Manage schema changes with goose, using plain SQL files in `api/migrations/`. Run them through `api/cmd/migrate`, which uses goose's core package and the existing `pgx` driver. Migrations run only when you ask (`make migrate`), not at server startup.

**Alternatives considered:**

- The goose command-line tool, pinned with the `tool` directive: this added 61 modules to `go.mod`, mostly drivers for databases we don't use.
- golang-migrate: similar to goose, but goose's Go package will also set up test databases.
- Atlas: more powerful, but more than this project needs.
- Running migrations at server startup: convenient, but schema changes would happen as a side effect of starting a server.

**Reasoning:** Our own command adds only 4 modules, and the code it uses is the same code that will set up test databases. Running migrations explicitly makes it clear when the schema changes.

## 2026-11-03: Defer deployment until the app is shared

**Decision:** Build and run Corkboard locally, with no hosting accounts, until there's a reason to show it to other people. Step 1 is complete without deploying. Keep `render.yaml` and the Pages Function in the repo, unused, so deploying later mostly means creating the accounts.

**Alternatives considered:**

- Deploy now on Render and Cloudflare Pages: free, but it adds accounts and online services before anyone needs to use the app.
- Remove the deployment config: a cleaner repo, but that work would need redoing later.

**Reasoning:** This is a learning project meant to grow slowly. Local Docker Postgres and CI on GitHub cover everything needed for now. Later steps that assume online services, such as email for sign-in and the digest, will use local alternatives.

## 2026-11-02: Use the API proxy from step 1 instead of CORS

**Decision:** The frontend always calls `/api` on its own origin. Vite's dev server forwards those requests locally, and a Cloudflare Pages Function will do the same when deployed. This moves the proxy from step 2 to the step 1 deploy.

**Alternatives considered:**

- CORS middleware in Go for step 1, followed by the proxy in step 2: this adds code we would soon delete.

**Reasoning:** Login needs the proxy anyway. Using it now means there's one path from browser to API, and the API never has to allow cross-origin requests.

## 2026-11-02: Pin TypeScript to 5.9.3

**Decision:** Pin the web project to TypeScript 5.9.3 instead of the template's TypeScript 7. Revisit this when `openapi-typescript` supports version 7.

**Alternatives considered:**

- Keep TypeScript 7 and run the generator through a pinned `npx` command: its dependencies stay outside the lockfile, and it downloads on each run.
- Keep TypeScript 7 and give the generator its own `package.json`: fully pinned, but a second set of dependencies for one small tool.
- Force the install with `--legacy-peer-deps`: this hides the conflict, and the generator may not work.

**Reasoning:** `openapi-typescript` 7.13.0 requires TypeScript 5, and npm also rejected an `overrides` workaround. One compiler and one lockfile keep generation repeatable for the drift check, and TypeScript 7's speed barely matters in an app this small.

## 2026-11-02: Use openapi-fetch for API calls

**Decision:** Pair `openapi-typescript` with `openapi-fetch` for frontend API calls. Accept a moderate `minimatch` audit finding in the generator's dev dependencies for now.

**Alternatives considered:**

- Plain `fetch` with handwritten types: types could drift from the contract.
- `npm audit fix --force`: this downgrades `openapi-typescript` to version 6.

**Reasoning:** `openapi-typescript` only generates types. `openapi-fetch` comes from the same project, is small, and rejects wrong paths and fields at compile time. The audit finding doesn't reach production, and the generator only reads our own spec, so CI should audit production dependencies with `npm audit --omit=dev`.

## 2026-11-01: Start with a small golangci-lint configuration

**Decision:** Use golangci-lint's standard linters plus `gosec` and `gofmt`, excluding generated code. The local version is 2.12.2, and CI will pin the same version.

**Alternatives considered:**

- Enable most linters: thorough, but noisy results tend to get ignored.
- Use only `go vet`: it misses ignored errors, dead code, and common security mistakes.

**Reasoning:** A short list keeps every finding worth acting on, and linters can be added as the code grows.

## 2026-10-31: Generate strict server handlers with oapi-codegen

**Decision:** Generate `std-http-server` code in strict mode, and pin the generator with the `tool` directive in `api/go.mod`.

**Alternatives considered:**

- Non-strict handlers using `http.ResponseWriter`: closer to raw `net/http`, but handlers can return a response the spec doesn't allow.
- Installing the generator globally: versions can differ between machines and CI.

**Reasoning:** Strict handlers return typed responses, so the compiler helps keep the server aligned with the contract. A pinned tool means local and CI generation produce the same output, which the drift check depends on.

## 2026-10-31: Use default hosting URLs, with a same-origin API proxy planned for auth

**Decision:** Do not buy a custom domain. Serve the frontend from `*.pages.dev` and the API from `*.onrender.com`. In step 2, route `/api` through a Cloudflare Pages Function so the browser sees one origin.

**Alternatives considered:**

- Custom domain with `api.` subdomain: the simplest cookie setup, but an ongoing cost that a pet project doesn't justify.
- Cross-site cookies with `SameSite=None`: third-party cookies are increasingly blocked, and this weakens the protection we want for an app used by teenagers.

**Reasoning:** `pages.dev` and `onrender.com` are separate sites, so browsers would treat an API session cookie as third-party. A proxy keeps cookies first-party without a domain, at the cost of one more moving part.

## 2026-10-31: Defer database wiring to step 2

**Decision:** Step 1 provides local Postgres through Docker Compose, but the API does not connect to it. Migrations, sqlc, and managed Postgres arrive with the first real table.

**Alternatives considered:**

- Connect now and report database status from the health check: this adds a dependency with nothing using it yet.
- Add migration tooling now: there is nothing to migrate.

**Reasoning:** The build plan asks us not to scaffold future features. This departs slightly from `CLAUDE.md`, which expects migration commands to be documented in step 1. They will be documented in step 2 instead.

## 2026-10-31: Pin Node to the current LTS line

**Decision:** Pin Node 26 in `.nvmrc`, used locally through fnm and in CI.

**Alternatives considered:**

- Stay on Node 25: odd-numbered releases never become LTS.
- Node 24: also LTS, but it reaches end of life sooner.

**Reasoning:** One pinned version keeps local and CI behavior the same, and an LTS line receives long-term security fixes.

## 2026-10-31: Use Docker Desktop for local containers

**Decision:** Use Docker Desktop to run local Postgres and, later, testcontainers.

**Alternatives considered:**

- Colima or OrbStack: workable, but testcontainers can need extra configuration on Colima, and OrbStack adds a commercial-licensing question.

**Reasoning:** It works with testcontainers out of the box and is free for personal use.
