# Makefile for Auth Service Docker Compose

.PHONY: help dev prod build up down logs clean backup restore shell

# Default environment is development
ENV ?= dev

# Default target
help: ## Show this help message
	@echo "Auth Service Docker Compose Management"
	@echo ""
	@echo "Usage: make [target] [ENV=dev|prod]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Environment setup
setup: ## Setup environment file
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example"; \
		cp .env.example .env; \
		echo "Please edit .env file with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

# Development environment
dev: setup ## Start development environment with tools
	@echo "Starting development environment..."
	@docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
	@make wait-for-services
	@make show-services

# Production environment  
prod: setup ## Start production environment
	@echo "Starting production environment..."
	@docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@make wait-for-services
	@make show-services

# Basic commands
build: ## Build all images
	@echo "Building images..."
	@if [ "$(ENV)" = "prod" ]; then \
		docker-compose -f docker-compose.yml -f docker-compose.prod.yml build; \
	else \
		docker-compose -f docker-compose.yml -f docker-compose.dev.yml build; \
	fi

up: setup ## Start services
	@echo "Starting services..."
	@if [ "$(ENV)" = "prod" ]; then \
		docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d; \
	else \
		docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d; \
	fi
	@make wait-for-services

down: ## Stop services
	@echo "Stopping services..."
	@docker-compose down

# Utility commands
logs: ## Show logs (usage: make logs [SERVICE=service_name])
	@if [ -n "$(SERVICE)" ]; then \
		docker-compose logs -f $(SERVICE); \
	else \
		docker-compose logs -f; \
	fi

status: ## Show service status
	@echo "Service Status:"
	@docker-compose ps
	@echo ""
	@make show-health

restart: ## Restart service (usage: make restart SERVICE=service_name)
	@if [ -z "$(SERVICE)" ]; then \
		echo "Usage: make restart SERVICE=service_name"; \
		exit 1; \
	fi
	@echo "Restarting $(SERVICE)..."
	@docker-compose restart $(SERVICE)

# Database operations
db-shell: ## Open PostgreSQL shell
	@docker-compose exec postgres psql -U postgres -d garmin_db

redis-shell: ## Open Redis shell
	@docker-compose exec redis redis-cli

backup: ## Backup database
	@echo "Creating database backup..."
	@mkdir -p backups
	@docker-compose exec postgres pg_dump -U postgres garmin_db | gzip > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql.gz
	@echo "Backup created in backups/ directory"

restore: ## Restore database (usage: make restore FILE=backup_file.sql.gz)
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make restore FILE=backup_file.sql.gz"; \
		exit 1; \
	fi
	@if [ ! -f "$(FILE)" ]; then \
		echo "Backup file not found: $(FILE)"; \
		exit 1; \
	fi
	@echo "Restoring database from $(FILE)..."
	@gunzip -c $(FILE) | docker-compose exec -T postgres psql -U postgres -d garmin_db
	@echo "Database restored successfully"

# Application operations
app-shell: ## Open shell in auth service container
	@docker-compose exec auth-service /bin/sh

app-build: ## Build only the auth service
	@echo "Building auth service..."
	@docker-compose build auth-service

app-logs: ## Show auth service logs
	@docker-compose logs -f auth-service

# Development tools
pgadmin: ## Open PgAdmin in browser
	@echo "Opening PgAdmin at http://localhost:8080"
	@echo "Login: admin@example.com / admin123"
	@open http://localhost:8080 2>/dev/null || xdg-open http://localhost:8080 2>/dev/null || echo "Please open http://localhost:8080 manually"

redis-commander: ## Open Redis Commander in browser
	@echo "Opening Redis Commander at http://localhost:8081"
	@open http://localhost:8081 2>/dev/null || xdg-open http://localhost:8081 2>/dev/null || echo "Please open http://localhost:8081 manually"

swagger: ## Open Swagger UI in browser
	@echo "Opening Swagger UI at http://localhost:5051/swagger/"
	@open http://localhost:5051/swagger/ 2>/dev/null || xdg-open http://localhost:5051/swagger/ 2>/dev/null || echo "Please open http://localhost:5051/swagger/ manually"

# Monitoring
monitoring: ## Start with monitoring tools
	@echo "Starting with monitoring tools..."
	@docker-compose --profile monitoring up -d
	@make wait-for-services
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000 (admin/admin123)"

# Cleanup operations
clean: ## Remove containers and volumes
	@echo "WARNING: This will remove all containers and volumes!"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@docker-compose down -v
	@docker system prune -f
	@echo "Cleanup completed"

clean-images: ## Remove all project images
	@echo "Removing project images..."
	@docker-compose down --rmi all
	@docker image prune -f

# Internal helper targets
wait-for-services:
	@echo "Waiting for services to be healthy..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker-compose ps postgres | grep -q "healthy" && \
		   docker-compose ps redis | grep -q "healthy"; then \
			echo "Services are healthy"; \
			break; \
		fi; \
		echo "Waiting for services... ($$timeout seconds remaining)"; \
		sleep 5; \
		timeout=$$((timeout-5)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "Services failed to become healthy"; \
		exit 1; \
	fi

show-services:
	@echo ""
	@echo "Services are running:"
	@echo "  Application:     http://localhost:5051"
	@echo "  Swagger UI:      http://localhost:5051/swagger/"
	@echo "  PostgreSQL:      localhost:5432"
	@echo "  Redis:           localhost:6379"
	@if docker-compose ps pgadmin | grep -q "Up"; then \
		echo "  PgAdmin:         http://localhost:8080"; \
	fi
	@if docker-compose ps redis-commander | grep -q "Up"; then \
		echo "  Redis Commander: http://localhost:8081"; \
	fi
	@if docker-compose ps grafana | grep -q "Up"; then \
		echo "  Grafana:         http://localhost:3000"; \
	fi
	@if docker-compose ps prometheus | grep -q "Up"; then \
		echo "  Prometheus:      http://localhost:9090"; \
	fi
	@echo ""

show-health:
	@echo "Health Status:"
	@for service in auth-service postgres redis; do \
		if docker-compose ps $$service | grep -q "healthy"; then \
			echo "  ✓ $$service"; \
		elif docker-compose ps $$service | grep -q "Up"; then \
			echo "  ⚠ $$service (no health check)"; \
		else \
			echo "  ✗ $$service"; \
		fi; \
	done

# Testing
test: ## Run tests in container
	@echo "Running tests..."
	@docker-compose exec auth-service go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@docker-compose exec auth-service go test -v -coverprofile=coverage.out ./...
	@docker-compose exec auth-service go tool cover -html=coverage.out -o coverage.html

# Quick development workflow
quick-restart: ## Quick restart of auth service
	@docker-compose restart auth-service
	@echo "Auth service restarted"

rebuild: ## Rebuild and restart auth service
	@echo "Rebuilding auth service..."
	@docker-compose build auth-service
	@docker-compose up -d auth-service
	@echo "Auth service rebuilt and restarted"