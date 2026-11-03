GENERATED_FILES := api/internal/api/api.gen.go web/src/api/schema.gen.ts

.DEFAULT_GOAL := help

.PHONY: help db-up db-down dev-api dev-web generate test lint typecheck build check-generated check

help: ## List available commands
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  %-16s %s\n", $$1, $$2}'

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

test: ## Run Go and web tests
	cd api && go test ./...
	npm --prefix web test

lint: ## Run go vet, golangci-lint, and Oxlint
	cd api && go vet ./...
	cd api && golangci-lint run ./...
	npm --prefix web run --silent lint

typecheck: ## Run the TypeScript check
	npm --prefix web run --silent typecheck

build: ## Build the API and web app
	cd api && go build ./...
	npm --prefix web run --silent build

check-generated: generate ## Fail if generated code is out of date with the spec
	@if [ -n "$$(git status --porcelain -- $(GENERATED_FILES))" ]; then \
		echo "Generated code is out of date. Run 'make generate' and commit the result."; \
		git status --short -- $(GENERATED_FILES); \
		exit 1; \
	fi

check: lint typecheck test build check-generated ## Run everything CI runs
