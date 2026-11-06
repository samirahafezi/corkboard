GENERATED_FILES := api/internal/api/api.gen.go web/src/api/schema.gen.ts

# Load local settings, such as DATABASE_URL, from .env when it exists, and
# pass them to every command. CI has no .env and sets what it needs itself.
ifneq (,$(wildcard .env))
include .env
export
endif

.DEFAULT_GOAL := help

.PHONY: help db-up db-down dev-api dev-web generate test test-api test-web lint lint-api lint-web typecheck build build-api build-web audit-web check-generated check

help: ## List available commands
	@grep -E '^[a-z-]+:.*## ' $(firstword $(MAKEFILE_LIST)) | awk -F ':.*## ' '{printf "  %-16s %s\n", $$1, $$2}'

db-up: ## Start local Postgres
	docker compose up -d --wait

db-down: ## Stop local Postgres and keep its data
	docker compose down

dev-api: ## Run the Go API server
	cd api && go run ./cmd/server

dev-web: ## Run the Vite dev server
	npm --prefix web run dev

generate: ## Regenerate Go and TypeScript code from the OpenAPI spec
	cd api && go generate ./...
	npm --prefix web run --silent generate

test: test-api test-web ## Run Go and web tests

test-api: ## Run Go tests
	cd api && go test ./...

test-web: ## Run web tests
	npm --prefix web test

lint: lint-api lint-web ## Run all linters

lint-api: ## Run go vet and golangci-lint
	cd api && go vet ./...
	cd api && golangci-lint run ./...

lint-web: ## Run Oxlint
	npm --prefix web run --silent lint

typecheck: ## Run the TypeScript check
	npm --prefix web run --silent typecheck

build: build-api build-web ## Build the API and web app

build-api: ## Build the Go API
	cd api && go build ./...

build-web: ## Build the web app
	npm --prefix web run --silent build

audit-web: ## Audit production web dependencies
	npm --prefix web audit --omit=dev

check-generated: generate ## Fail if generated code is out of date with the spec
	@if [ -n "$$(git status --porcelain -- $(GENERATED_FILES))" ]; then \
		echo "Generated code is out of date. Run 'make generate' and commit the result."; \
		git status --short -- $(GENERATED_FILES); \
		exit 1; \
	fi

check: lint typecheck test build audit-web check-generated ## Run everything CI runs
