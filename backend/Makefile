# Idea Collision Engine API - Makefile

.PHONY: build run test clean deps migrate docker-build docker-run help

# Variables
APP_NAME = idea-collision-engine-api
BINARY_NAME = main
MIGRATE_BINARY = migrate
GO_VERSION = 1.22

# Default target
help: ## Show available commands
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Dependencies
deps: ## Download Go dependencies
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

# Build
build: deps ## Build the application
	@echo "🔨 Building application..."
	go build -o $(BINARY_NAME) ./cmd/server
	go build -o $(MIGRATE_BINARY) ./cmd/migrate

# Run development server
run: build ## Run the development server
	@echo "🚀 Starting development server..."
	./$(BINARY_NAME)

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev: deps ## Run with live reload
	@if command -v air > /dev/null; then \
		echo "🔥 Starting development server with live reload..."; \
		air; \
	else \
		echo "❌ Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "   Or run: make run"; \
	fi

# Database operations
migrate: build ## Run database migrations
	@echo "🗄️  Running database migrations..."
	./$(MIGRATE_BINARY)

db-reset: ## Reset database (requires manual confirmation)
	@echo "⚠️  This will drop and recreate the database!"
	@echo "   Make sure DATABASE_URL is set correctly"
	@read -p "   Continue? [y/N]: " confirm && [ "$$confirm" = "y" ]
	@echo "🗑️  Resetting database..."
	psql $(DATABASE_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	$(MAKE) migrate

# Testing
test: deps ## Run tests
	@echo "🧪 Running tests..."
	go run ../test_runner.go

test-coverage: test ## Run tests with coverage report
	@echo "📊 Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "   Coverage report: coverage.html"

benchmark: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...

# Code quality
lint: ## Run linters
	@if command -v golangci-lint > /dev/null; then \
		echo "🔍 Running linters..."; \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not found. Install from: https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

# Docker operations
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):latest .

docker-run: docker-build ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

docker-compose-up: ## Start services with docker-compose
	@echo "🐳 Starting services with docker-compose..."
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	@echo "🐳 Stopping services with docker-compose..."
	docker-compose down

# Deployment
deploy-staging: ## Deploy to staging (Fly.io)
	@echo "🚀 Deploying to staging..."
	flyctl deploy --config fly.staging.toml

deploy-prod: ## Deploy to production (Fly.io)
	@echo "🚀 Deploying to production..."
	@echo "⚠️  This will deploy to PRODUCTION!"
	@read -p "   Continue? [y/N]: " confirm && [ "$$confirm" = "y" ]
	flyctl deploy --config fly.toml

# Utility commands
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -f $(BINARY_NAME) $(MIGRATE_BINARY)
	rm -f coverage.out coverage.html
	go clean

logs: ## View application logs (Fly.io)
	flyctl logs -a $(APP_NAME)

status: ## Check application status (Fly.io)
	flyctl status -a $(APP_NAME)

ssh: ## SSH into production container (Fly.io)
	flyctl ssh console -a $(APP_NAME)

# Environment setup
env-example: ## Create .env file from example
	@if [ ! -f .env ]; then \
		echo "📝 Creating .env file from example..."; \
		cp .env.example .env; \
		echo "   Please edit .env with your configuration"; \
	else \
		echo "❌ .env file already exists"; \
	fi

# Security
security-scan: ## Run security scan
	@if command -v gosec > /dev/null; then \
		echo "🔒 Running security scan..."; \
		gosec ./...; \
	else \
		echo "❌ gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Performance profiling
profile-cpu: build ## Run CPU profiling
	@echo "⚡ Running CPU profiling..."
	go tool pprof -http=:8081 http://localhost:8080/debug/pprof/profile

profile-mem: build ## Run memory profiling  
	@echo "⚡ Running memory profiling..."
	go tool pprof -http=:8081 http://localhost:8080/debug/pprof/heap

# Development setup
setup: deps env-example ## Setup development environment
	@echo "🛠️  Setting up development environment..."
	@echo "✅ Dependencies downloaded"
	@echo "✅ Environment file created"
	@echo ""
	@echo "Next steps:"
	@echo "1. Edit .env with your configuration"
	@echo "2. Set up PostgreSQL and Redis"
	@echo "3. Run 'make migrate' to setup database"
	@echo "4. Run 'make dev' to start development server"

# CI/CD helpers
ci-test: deps test lint vet ## Run all CI checks
	@echo "✅ All CI checks passed!"

# Database backup/restore (for local development)
db-backup: ## Backup database
	@echo "💾 Creating database backup..."
	pg_dump $(DATABASE_URL) > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Show environment info
info: ## Show environment information
	@echo "Environment Information:"
	@echo "  Go Version: $(shell go version)"
	@echo "  App Name: $(APP_NAME)"
	@echo "  Git Branch: $(shell git branch --show-current 2>/dev/null || echo 'unknown')"
	@echo "  Git Commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "  Build Time: $(shell date)"