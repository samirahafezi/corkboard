# Corkboard

**A calm, social-media-free announcement board for high school clubs.**

Club leads post announcements. Staff advisors approve them. Students follow the clubs they care about and see updates in a simple chronological feed, or in a weekly email digest if they'd rather never open an app at all.

> 🚧 **Status:** early development. See the [roadmap](#roadmap) for progress.

## Why

More and more schools share club updates and announcements on Instagram. That quietly pushes students into creating social media accounts just to find out when the chess club meets, and it leaves out students and families who choose not to be on social media.

Corkboard is the digital version of the corkboard outside the school office: one trustworthy place for school updates, with none of the engagement machinery.

## Principles

- **Calm by design.** No likes, view counts, streaks, or algorithmic feeds. The feed is chronological, and it ends.
- **Privacy first.** Built for teenagers, so it collects as little data as possible and never sells or shares it.
- **Works even if you never open it.** A weekly email digest delivers what you follow.
- **Accessible.** Targets WCAG 2.1 AA, with automated accessibility checks in CI.

## How it works

1. A **club lead** drafts an announcement, with an optional event date, time, and location.
2. The club's **staff advisor** reviews it, then publishes it or returns it with a note.
3. **Students** who follow the club see it in their feed and in their weekly digest.

## Tech stack

| Layer | Technology |
| --- | --- |
| Frontend | React, Vite, TypeScript, TanStack Query (PWA) |
| Backend | Go |
| API contract | OpenAPI, with generated Go server code and TypeScript client |
| Database | Postgres, sqlc |
| Auth | Passwordless magic-link sign-in |
| Testing | Go testing, testcontainers, Vitest, Playwright, axe |
| CI/CD | GitHub Actions |
| Hosting | Cloudflare Pages, Render |

## Engineering practices

This project is also an exercise in building software well:

- **Contract-first API.** The OpenAPI spec is the single source of truth; both server and client code are generated from it, and CI fails if they drift.
- **Quality gates.** Tests, linting, type checks, and accessibility checks block every merge.
- **Security by default.** Hashed single-use sign-in tokens, secure session cookies, and authorization tests proving each role can only do what it should.
- **AI-assisted, human-reviewed.** Built with AI pair programming, with every change reviewed and verified by tests. Notable decisions are recorded in [DECISIONS.md](DECISIONS.md).

## Roadmap

- [x] Foundation: repo structure, CI pipeline, API contract
- [ ] Passwordless sign-in and roles
- [ ] Club directory and following
- [ ] Announcements and advisor approval
- [ ] Student feed
- [ ] Weekly email digest
- [ ] Polish, demo data, and write-up

**Later ideas:** multiple schools, calendar feeds, push notifications, images.

## Getting started

Setup instructions are coming as part of the foundation step.

## License

To be decided.# Corkboard
