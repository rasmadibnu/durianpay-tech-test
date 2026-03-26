.PHONY: help install run-backend run-frontend run build test test-backend test-frontend test-frontend-typecheck clean

.DEFAULT_GOAL := help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install all dependencies
	cd backend && go mod tidy && go mod vendor
	cd frontend && npm install

run-backend: ## Run backend dev server
	cd backend && CGO_ENABLED=1 go run main.go

run-frontend: ## Run frontend dev server
	cd frontend && npm run dev

run: ## Run both services with docker-compose
	docker compose up --build

build: ## Build both services with docker-compose
	docker compose build

test-backend: ## Run backend tests
	cd backend && CGO_ENABLED=1 go test ./... -v

test-frontend: ## Run frontend unit tests
	cd frontend && npm test

test-frontend-typecheck: ## Run frontend TypeScript check
	cd frontend && npx tsc --noEmit

test: test-backend test-frontend ## Run all tests

clean: ## Clean build artifacts
	docker compose down -v
	rm -f backend/dashboard.db
	rm -rf frontend/dist
