# URL Shortener Makefile

# Variables
DOCKER_COMPOSE = docker-compose
BACKEND_DIR = backend
FRONTEND_DIR = frontend

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m # No Color

.PHONY: help build up down logs clean test backend frontend dev

# Default target
help: ## Show this help message
	@echo "URL Shortener - Available Commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Start development environment with hot reload
	@echo "$(YELLOW)Starting development environment...$(NC)"
	$(DOCKER_COMPOSE) --profile dev up -d
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo "Frontend (dev): http://localhost:3001"
	@echo "Backend API: http://localhost:8080"
	@echo "Database: postgresql://postgres:password@localhost:5432/urlshortener"

# Production
up: ## Start all services in production mode
	@echo "$(YELLOW)Starting production environment...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Production environment started!$(NC)"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:8080"
	@echo "Database: postgresql://postgres:password@localhost:5432/urlshortener"

# Build
build: ## Build all Docker images
	@echo "$(YELLOW)Building Docker images...$(NC)"
	$(DOCKER_COMPOSE) build
	@echo "$(GREEN)Build completed!$(NC)"

# Stop services
down: ## Stop all services
	@echo "$(YELLOW)Stopping services...$(NC)"
	$(DOCKER_COMPOSE) --profile dev --profile cache --profile tools down
	@echo "$(GREEN)Services stopped!$(NC)"

# Database
db: ## Start database and tools
	@echo "$(YELLOW)Starting database services...$(NC)"
	$(DOCKER_COMPOSE) --profile tools up -d postgres adminer
	@echo "$(GREEN)Database services started!$(NC)"
	@echo "Database: postgresql://postgres:password@localhost:5432/urlshortener"
	@echo "Adminer: http://localhost:8081"

# Cache
cache: ## Start with Redis cache
	@echo "$(YELLOW)Starting services with Redis cache...$(NC)"
	$(DOCKER_COMPOSE) --profile cache up -d
	@echo "$(GREEN)Services with cache started!$(NC)"

# Logs
logs: ## Show logs from all services
	$(DOCKER_COMPOSE) logs -f

logs-backend: ## Show backend logs
	$(DOCKER_COMPOSE) logs -f backend

logs-frontend: ## Show frontend logs
	$(DOCKER_COMPOSE) logs -f frontend

# Testing
test: ## Run tests
	@echo "$(YELLOW)Running backend tests...$(NC)"
	cd $(BACKEND_DIR) && go test ./tests/... -v
	@echo "$(YELLOW)Running frontend tests...$(NC)"
	cd $(FRONTEND_DIR) && npm test
	@echo "$(GREEN)Tests completed!$(NC)"

test-backend: ## Run backend tests only
	@echo "$(YELLOW)Running backend tests...$(NC)"
	cd $(BACKEND_DIR) && go test ./tests/... -v

test-frontend: ## Run frontend tests only
	@echo "$(YELLOW)Running frontend tests...$(NC)"
	cd $(FRONTEND_DIR) && npm test

# Backend specific
backend: ## Build and run backend only
	@echo "$(YELLOW)Starting backend...$(NC)"
	$(DOCKER_COMPOSE) up -d postgres backend
	@echo "$(GREEN)Backend started!$(NC)"

backend-logs: ## Show backend logs
	$(DOCKER_COMPOSE) logs -f backend

# Frontend specific
frontend: ## Build and run frontend only
	@echo "$(YELLOW)Starting frontend...$(NC)"
	$(DOCKER_COMPOSE) up -d frontend
	@echo "$(GREEN)Frontend started!$(NC)"

frontend-logs: ## Show frontend logs
	$(DOCKER_COMPOSE) logs -f frontend

# Cleanup
clean: ## Remove all containers, volumes, and images
	@echo "$(RED)Cleaning up Docker resources...$(NC)"
	$(DOCKER_COMPOSE) --profile dev --profile cache --profile tools down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)Cleanup completed!$(NC)"

clean-volumes: ## Remove all volumes (WARNING: This will delete all data)
	@echo "$(RED)Removing all volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)Volumes removed!$(NC)"

# Health check
health: ## Check service health
	@echo "$(YELLOW)Checking service health...$(NC)"
	@curl -f http://localhost:8080/health || echo "$(RED)Backend not responding$(NC)"
	@curl -f http://localhost:3000 || echo "$(RED)Frontend not responding$(NC)"

# Database operations
db-migrate: ## Run database migrations
	@echo "$(YELLOW)Running database migrations...$(NC)"
	$(DOCKER_COMPOSE) exec backend ./main migrate

db-seed: ## Seed database with sample data
	@echo "$(YELLOW)Seeding database...$(NC)"
	$(DOCKER_COMPOSE) exec backend ./main seed

# Install dependencies
install: ## Install dependencies for both frontend and backend
	@echo "$(YELLOW)Installing backend dependencies...$(NC)"
	cd $(BACKEND_DIR) && go mod download
	@echo "$(YELLOW)Installing frontend dependencies...$(NC)"
	cd $(FRONTEND_DIR) && npm install
	@echo "$(GREEN)Dependencies installed!$(NC)"

# Linting
lint: ## Run linters
	@echo "$(YELLOW)Running backend linter...$(NC)"
	cd $(BACKEND_DIR) && go fmt ./...
	@echo "$(YELLOW)Running frontend linter...$(NC)"
	cd $(FRONTEND_DIR) && npm run lint
	@echo "$(GREEN)Linting completed!$(NC)"

# Quick start
quick-start: build up ## Build and start everything
	@echo "$(GREEN)Quick start completed!$(NC)"

# Show status
status: ## Show running containers
	@echo "$(YELLOW)Container Status:$(NC)"
	$(DOCKER_COMPOSE) ps