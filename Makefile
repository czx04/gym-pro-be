.PHONY: help run build test clean docker-up docker-down migrate-up migrate-down migrate-create swagger-gen fmt lint

# Variables
APP_NAME=gym-pro-api
DOCKER_COMPOSE=docker-compose
GO=go
GOBIN=$(shell go env GOPATH)/bin
MIGRATE=$(GOBIN)/migrate
SWAG=$(GOBIN)/swag

# Help command
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
run: ## Run the application locally
	$(GO) run cmd/api/main.go

build: ## Build the application
	$(GO) build -o bin/$(APP_NAME) cmd/api/main.go

test: ## Run tests
	$(GO) test -v -race -cover ./...

test-coverage: ## Run tests with coverage
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Code quality
fmt: ## Format code
	$(GO) fmt ./...
	gofmt -s -w .

lint: ## Run linter
	golangci-lint run --timeout=5m

# Dependencies
deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

deps-upgrade: ## Upgrade dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

# Docker
docker-up: ## Start Docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-build: ## Build Docker image
	docker build -t $(APP_NAME):latest -f docker/Dockerfile .

docker-rebuild: ## Rebuild and restart Docker containers
	$(DOCKER_COMPOSE) up -d --build

# Database migrations
migrate-install: ## Install golang-migrate tool
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Run database migrations up
	$(MIGRATE) -path migrations -database "postgresql://gymadmin:secret123@localhost:5432/gym_pro_db?sslmode=disable" up

migrate-down: ## Roll back database migrations
	$(MIGRATE) -path migrations -database "postgresql://gymadmin:secret123@localhost:5432/gym_pro_db?sslmode=disable" down

migrate-drop: ## Drop all database tables
	$(MIGRATE) -path migrations -database "postgresql://gymadmin:secret123@localhost:5432/gym_pro_db?sslmode=disable" drop -f

migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	$(MIGRATE) create -ext sql -dir migrations -seq $(name)

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	$(MIGRATE) -path migrations -database "postgresql://gymadmin:secret123@localhost:5432/gym_pro_db?sslmode=disable" force $(version)

# Swagger
swagger-install: ## Install swag tool
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

swagger-gen: ## Generate Swagger documentation
	$(SWAG) init -g cmd/api/main.go -o docs --parseDependency --parseInternal

swagger-fmt: ## Format Swagger comments
	$(SWAG) fmt

# Database
db-create: ## Create database
	docker exec -it gym-pro-postgres psql -U gymadmin -c "CREATE DATABASE gym_pro_db;"

db-drop: ## Drop database
	docker exec -it gym-pro-postgres psql -U gymadmin -c "DROP DATABASE IF EXISTS gym_pro_db;"

db-reset: db-drop db-create migrate-up ## Reset database (drop, create, migrate)

db-psql: ## Connect to database with psql
	docker exec -it gym-pro-postgres psql -U gymadmin -d gym_pro_db

# Development workflow
dev-setup: deps migrate-install swagger-install ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "Run 'make docker-up' to start PostgreSQL"
	@echo "Run 'make migrate-up' to run migrations"
	@echo "Run 'make run' to start the API server"

dev-start: docker-up ## Start all development services
	@echo "Waiting for database to be ready..."
	@sleep 3
	@make migrate-up
	@echo "Starting API server..."
	@make run

# Production
prod-build: ## Build for production
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-w -s" -o bin/$(APP_NAME) cmd/api/main.go

# Git
git-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@echo "#!/bin/sh\nmake fmt\nmake lint" > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed!"
