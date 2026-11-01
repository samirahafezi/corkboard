# Decisions

A short log of meaningful technical choices, newest first. The initial stack is listed in `CLAUDE.md`.

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
