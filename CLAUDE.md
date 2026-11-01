# CLAUDE.md

Guidance for Claude Code when working in this repository.

## What this project is

Corkboard is a calm, social-media-free announcement app for high school clubs. Club leads post announcements, staff advisors approve them, and students who follow a club see them in a chronological feed and a weekly email digest.

It exists because schools increasingly push club and school updates through Instagram, which forces teens to create social media accounts just to stay informed. Corkboard is the digital version of the corkboard outside the school office: useful, quiet, and not designed to be addictive.

This is a personal side project built for learning and as a portfolio piece. The owner is an experienced engineering manager moving toward an IC role, and is using this project to deepen Go skills and demonstrate strong engineering practices: quality gates, testing, security, and a thoughtful AI-assisted workflow.

## How to work with me

These rules matter more than speed.

- **One step at a time.** Work only on the current step of the build plan below. Do not jump ahead or scaffold future features.
- **Plan before code.** At the start of each task, explain your plan and wait for my go-ahead before writing code.
- **Teach the Go.** I'm learning idiomatic Go. When you introduce a Go pattern, idiom, or standard library feature, briefly explain what it is and why it's the idiomatic choice. When a piece is a good learning exercise, offer to let me write it myself and review my version.
- **Small, focused commits.** One logical change per commit, with clear messages (Conventional Commits style: `feat:`, `fix:`, `test:`, `chore:`, `docs:`, `ci:`).
- **Never push, merge, or open PRs without asking.**
- **Checks pass locally first.** Run the same checks CI runs before calling a task done.
- **Log decisions.** When we make a meaningful technical choice, add a short entry to `DECISIONS.md` (date, decision, alternatives considered, reasoning).
- **Flag uncertainty.** If you're unsure about something, say so rather than guessing. If you notice a problem outside the current task, mention it instead of fixing it silently.
- **No secrets in the repo.** Configuration comes from environment variables. `.env` files stay in `.gitignore`. Provide a `.env.example` with placeholder values instead.

## Product principles

These shape design decisions, not just copy:

- **Calm by design.** No likes, view counts, streaks, infinite scroll, or engagement algorithms. The feed is chronological and ends.
- **Privacy first.** The users are teenagers. Collect the minimum data needed, never sell or share it, and prefer deleting old data over keeping it.
- **Works for people who never open the app.** The email digest is a first-class feature, not an afterthought.
- **Accessible.** Meet WCAG 2.1 AA. Accessibility checks run in CI.

## Stack

| Layer | Choice |
| --- | --- |
| Frontend | React + Vite + TypeScript, TanStack Query, installable as a PWA |
| Backend | Go, standard library `net/http` routing (or chi if needed) |
| API contract | OpenAPI spec in `openapi/`; Go server code generated with oapi-codegen, TypeScript client with openapi-typescript |
| Database | Postgres, queries with sqlc |
| Auth | Our own magic-link authentication in Go, with secure session cookies |
| Authorization | Role checks in Go middleware |
| Email | Resend |
| Testing | Go `testing` + testcontainers (real Postgres), Vitest, Playwright, axe |
| CI | GitHub Actions |
| Local dev | Postgres in Docker |
| Hosting | Cloudflare Pages (frontend), Render (API), separate managed Postgres |

## Repo layout

```
corkboard/
  api/          Go backend
  web/          React frontend
  openapi/      shared API contract (source of truth)
  .github/      CI workflows
  CLAUDE.md
  DECISIONS.md
  README.md
```

## API contract rules

- `openapi/openapi.yaml` is the single source of truth for the API.
- Change the spec first, then regenerate code. Never hand-edit generated files.
- CI fails if generated code is out of date with the spec.

## Security expectations

- Magic-link tokens: cryptographically random, stored hashed, short expiry, single use.
- Session cookies: `HttpOnly`, `Secure`, `SameSite`.
- Every endpoint enforces authorization in middleware, and has tests proving that roles without access are rejected.
- Parameterized queries only (sqlc handles this). Validate all input at the API boundary.

## Roles

| Role | Can do |
| --- | --- |
| Student | Browse clubs, follow and unfollow, read published announcements, manage digest settings |
| Club lead | Everything a student can, plus draft and submit announcements for their own club |
| Staff advisor | Approve or return announcements for clubs they advise |
| School admin | Manage clubs, assign club leads and advisors |

## Announcement states

`draft` → `pending_review` → `published`, or `pending_review` → `returned` (with a note) → back to `draft`. Every valid and invalid transition is tested.

## v1 scope

**In:** one school, the four roles above, club directory with following, announcements with optional event date/time/location, the approval workflow, a chronological feed, a weekly email digest.

**Out (do not build yet):** comments, DMs, images, multiple schools, push notifications, native apps, calendar (ICS) feeds, Postgres row-level security.

## Build plan

Current step: **1. Foundation**

1. **Foundation.** Monorepo structure, Go health-check endpoint, React page that calls it, OpenAPI spec with code generation, local Postgres in Docker. CI runs Go tests, `go vet`, linting, TypeScript checks, and an OpenAPI drift check, all blocking merges. Both halves deployed.
2. **Auth and roles.** Magic-link sign-in, sessions, user roles, authorization middleware, and tests proving each role can't do what it shouldn't.
3. **Clubs and following.** Club directory, club pages, follow and unfollow.
4. **Announcements and approval.** The state machine above, with tests for every transition.
5. **The feed.** Chronological feed of published posts from followed clubs.
6. **Weekly digest.** Scheduled Go job emailing each student their followed clubs' posts.
7. **Polish and write-up.** Demo seed data, architecture notes in the README, accessibility pass, and a write-up of the AI-assisted workflow.

Update "Current step" above when a step is complete.

## Commands

To be filled in during step 1 (dev server, tests, linting, code generation, database migrations).
