# =============================================================================
# project-box-go-fiber-qwik-city Makefile
# =============================================================================
# Usage:
#   make help          - Show this help
#   make dev           - Start development environment
#   make build         - Build for production
#   make test          - Run all tests
#   make deploy        - Deploy to production
# =============================================================================

.PHONY: help dev dev-backend dev-frontend build build-backend build-frontend \
        test test-backend test-frontend lint lint-backend lint-frontend \
        format docker docker-build docker-up docker-down docker-logs \
        clean install deploy

# Default target
.DEFAULT_GOAL := help

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

# =============================================================================
# HELP
# =============================================================================

help: ## Show this help message
	@echo ""
	@echo "$(CYAN)project-box-go-fiber-qwik-city$(RESET)"
	@echo "$(CYAN)==============================$(RESET)"
	@echo ""
	@echo "$(GREEN)Available commands:$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# DEVELOPMENT
# =============================================================================

install: ## Install all dependencies
	@echo "$(CYAN)Installing backend dependencies...$(RESET)"
	cd backend-go-fiber && go mod download
	@echo "$(CYAN)Installing frontend dependencies...$(RESET)"
	cd frontend-qwik-city && npm install
	@echo "$(GREEN)All dependencies installed$(RESET)"

dev: ## Start both backend and frontend in development mode (requires tmux or run in separate terminals)
	@echo "$(YELLOW)Start backend and frontend in separate terminals:$(RESET)"
	@echo "  Terminal 1: make dev-backend"
	@echo "  Terminal 2: make dev-frontend"

dev-backend: ## Start backend in development mode
	@echo "$(CYAN)Starting Go Fiber backend...$(RESET)"
	cd backend-go-fiber && go run cmd/server/main.go

dev-frontend: ## Start frontend in development mode
	@echo "$(CYAN)Starting Qwik City frontend...$(RESET)"
	cd frontend-qwik-city && npm run dev

# =============================================================================
# BUILD
# =============================================================================

build: build-backend build-frontend ## Build both backend and frontend for production
	@echo "$(GREEN)Production build complete$(RESET)"

build-backend: ## Build backend binary
	@echo "$(CYAN)Building Go backend...$(RESET)"
	cd backend-go-fiber && CGO_ENABLED=1 go build -o bin/server ./cmd/server
	@echo "$(GREEN)Backend built: backend-go-fiber/bin/server$(RESET)"

build-frontend: ## Build frontend for production
	@echo "$(CYAN)Building Qwik City frontend...$(RESET)"
	cd frontend-qwik-city && npm run build
	@echo "$(GREEN)Frontend built: frontend-qwik-city/dist/$(RESET)"

# =============================================================================
# TESTING
# =============================================================================

test: test-backend test-frontend ## Run all tests
	@echo "$(GREEN)All tests passed$(RESET)"

test-backend: ## Run backend tests
	@echo "$(CYAN)Running Go tests...$(RESET)"
	cd backend-go-fiber && go test -v ./...

test-backend-coverage: ## Run backend tests with coverage
	@echo "$(CYAN)Running Go tests with coverage...$(RESET)"
	cd backend-go-fiber && go test -v -coverprofile=coverage.out ./...
	cd backend-go-fiber && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: backend-go-fiber/coverage.html$(RESET)"

test-frontend: ## Run frontend checks
	@echo "$(CYAN)Running Qwik City build check...$(RESET)"
	cd frontend-qwik-city && npm run build

# =============================================================================
# LINTING & FORMATTING
# =============================================================================

lint: lint-backend lint-frontend ## Lint both backend and frontend
	@echo "$(GREEN)Linting complete$(RESET)"

lint-backend: ## Lint Go code (requires golangci-lint)
	@echo "$(CYAN)Linting Go code...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd backend-go-fiber && golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed, running go vet instead$(RESET)"; \
		cd backend-go-fiber && go vet ./...; \
	fi

lint-frontend: ## Lint frontend code
	@echo "$(CYAN)Linting frontend code...$(RESET)"
	cd frontend-qwik-city && npm run lint

format: format-backend format-frontend ## Format all code
	@echo "$(GREEN)Formatting complete$(RESET)"

format-backend: ## Format Go code
	@echo "$(CYAN)Formatting Go code...$(RESET)"
	cd backend-go-fiber && go fmt ./...

format-frontend: ## Format frontend code
	@echo "$(CYAN)Formatting frontend code...$(RESET)"
	cd frontend-qwik-city && npm run format

# =============================================================================
# DOCKER
# =============================================================================

docker: docker-build docker-up ## Build and start Docker containers

docker-build: ## Build Docker images
	@echo "$(CYAN)Building Docker images...$(RESET)"
	docker-compose build
	@echo "$(GREEN)Docker images built$(RESET)"

docker-up: ## Start Docker containers
	@echo "$(CYAN)Starting Docker containers...$(RESET)"
	@if [ ! -f .env ]; then \
		echo "$(RED)Error: .env file not found. Copy .env.example to .env and configure it.$(RESET)"; \
		exit 1; \
	fi
	docker-compose up -d
	@echo "$(GREEN)Containers started$(RESET)"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Backend:  http://localhost:3001"

docker-down: ## Stop Docker containers
	@echo "$(CYAN)Stopping Docker containers...$(RESET)"
	docker-compose down
	@echo "$(GREEN)Containers stopped$(RESET)"

docker-logs: ## Show Docker container logs
	docker-compose logs -f

docker-postgres: ## Start with PostgreSQL profile
	@echo "$(CYAN)Starting with PostgreSQL...$(RESET)"
	docker-compose --profile postgres up -d
	@echo "$(GREEN)Containers started with PostgreSQL$(RESET)"

# =============================================================================
# UTILITIES
# =============================================================================

clean: ## Clean build artifacts
	@echo "$(CYAN)Cleaning build artifacts...$(RESET)"
	rm -rf backend-go-fiber/bin
	rm -rf backend-go-fiber/coverage.out
	rm -rf backend-go-fiber/coverage.html
	rm -rf frontend-qwik-city/dist
	rm -rf frontend-qwik-city/.qwik
	rm -rf frontend-qwik-city/node_modules/.vite
	rm -rf frontend-qwik-city/tmp
	rm -rf frontend-qwik-city/server
	@echo "$(GREEN)Clean complete$(RESET)"

env-setup: ## Create .env file from example
	@if [ -f .env ]; then \
		echo "$(YELLOW).env already exists. Skipping.$(RESET)"; \
	else \
		cp .env.example .env; \
		echo "$(GREEN)Created .env from .env.example$(RESET)"; \
		echo "$(YELLOW)Remember to update JWT_SECRET before production!$(RESET)"; \
	fi

generate-secret: ## Generate a random JWT secret
	@echo "$(CYAN)Generated JWT_SECRET:$(RESET)"
	@openssl rand -base64 32

# =============================================================================
# DATABASE
# =============================================================================

db-reset: ## Reset SQLite database (WARNING: destroys data)
	@echo "$(RED)This will delete all data. Press Ctrl+C to cancel.$(RESET)"
	@sleep 3
	rm -f data/db/sqlite/app.db
	@echo "$(GREEN)Database reset$(RESET)"

seed: ## Seed database with test data
	@echo "$(CYAN)Seeding database...$(RESET)"
	cd backend-go-fiber && go run cmd/seed/main.go
	@echo "$(GREEN)Database seeded$(RESET)"

db-fresh: db-reset seed ## Reset and seed database

# =============================================================================
# DEPLOY
# =============================================================================

deploy: ## Deploy to production (customize for your setup)
	@echo "$(CYAN)Deploying to production...$(RESET)"
	@echo "$(YELLOW)Customize this target for your deployment setup$(RESET)"
	@echo ""
	@echo "Example deployment steps:"
	@echo "  1. Build: make build"
	@echo "  2. Push to registry: docker-compose push"
	@echo "  3. Deploy on server: docker-compose pull && docker-compose up -d"

deploy-check: ## Pre-deployment checklist
	@echo "$(CYAN)Pre-deployment checklist:$(RESET)"
	@echo ""
	@echo "  [ ] NODE_ENV=production in .env"
	@echo "  [ ] JWT_SECRET is 32+ characters and random"
	@echo "  [ ] CORS_ORIGINS set to production domain"
	@echo "  [ ] DATABASE_URL set to production database"
	@echo "  [ ] SSL/TLS configured (nginx or load balancer)"
	@echo "  [ ] Backup strategy in place"
	@echo "  [ ] Monitoring configured"
	@echo ""
