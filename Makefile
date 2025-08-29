# Variables
BINARY_NAME=sample-grpc-server
DOCKER_COMPOSE=docker-compose
GO=go
GOMOD=$(GO) mod
GOFMT=gofmt
GOLINT=golangci-lint
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOGET=$(GO) get

# Paths
CMD_PATH=./cmd/server
SEED_PATH=./cmd/seed
PROTO_PATH=./api/proto
GEN_PATH=./pkg/generated

# Build tags
BUILD_TAGS=
TEST_TAGS=

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help
help: ## Display this help screen
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
	@echo "$(GREEN)Building binary...$(NC)"
	$(GOBUILD) -o $(BINARY_NAME) -v $(CMD_PATH)
	@echo "$(GREEN)Build complete!$(NC)"

.PHONY: run
run: ## Run the application
	@echo "$(GREEN)Running application...$(NC)"
	$(GO) run $(CMD_PATH)/main.go

.PHONY: clean
clean: ## Remove build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(GEN_PATH)
	@echo "$(GREEN)Clean complete!$(NC)"

.PHONY: test
test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)Tests complete!$(NC)"

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

.PHONY: deps
deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies downloaded!$(NC)"

.PHONY: fmt
fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	@echo "$(GREEN)Formatting complete!$(NC)"

.PHONY: lint
lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	@if ! which golangci-lint > /dev/null; then \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run --timeout 5m
	@echo "$(GREEN)Linting complete!$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)Vet complete!$(NC)"

.PHONY: migrate
migrate: ## Run database migrations
	@echo "$(GREEN)Running database migrations...$(NC)"
	@$(GO) run $(CMD_PATH)/../migrate/main.go up
	@echo "$(GREEN)Migrations complete!$(NC)"

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)Rolling back database migrations...$(NC)"
	@$(GO) run $(CMD_PATH)/../migrate/main.go down
	@echo "$(GREEN)Rollback complete!$(NC)"

.PHONY: seed
seed: ## Seed the database with sample data
	@echo "$(GREEN)Seeding database...$(NC)"
	@$(GO) run $(SEED_PATH)/main.go
	@echo "$(GREEN)Database seeded successfully!$(NC)"

.PHONY: seed-clean
seed-clean: ## Clean and reseed the database
	@echo "$(YELLOW)Cleaning and reseeding database...$(NC)"
	@$(GO) run $(SEED_PATH)/main.go --clean
	@echo "$(GREEN)Database cleaned and reseeded!$(NC)"

.PHONY: proto
proto: ## Generate code from proto files
	@echo "$(GREEN)Generating proto files...$(NC)"
	@bash scripts/generate-proto.sh
	@echo "$(GREEN)Proto generation complete!$(NC)"

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):latest .
	@echo "$(GREEN)Docker build complete!$(NC)"

.PHONY: docker-up
docker-up: ## Start Docker containers
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Containers started!$(NC)"

.PHONY: docker-down
docker-down: ## Stop Docker containers
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Containers stopped!$(NC)"

.PHONY: docker-logs
docker-logs: ## View Docker logs
	$(DOCKER_COMPOSE) logs -f

.PHONY: docker-dev
docker-dev: ## Start development Docker containers with hot reload
	@echo "$(GREEN)Starting development containers...$(NC)"
	$(DOCKER_COMPOSE) --profile dev up -d
	@echo "$(GREEN)Development containers started!$(NC)"
	@echo "$(YELLOW)Server running with hot reload on port 50051$(NC)"
	@echo "$(YELLOW)Adminer available at http://localhost:8080$(NC)"

.PHONY: docker-dev-build
docker-dev-build: ## Build development Docker image
	@echo "$(GREEN)Building development Docker image...$(NC)"
	$(DOCKER_COMPOSE) --profile dev build
	@echo "$(GREEN)Development build complete!$(NC)"

.PHONY: docker-dev-logs
docker-dev-logs: ## View development server logs
	$(DOCKER_COMPOSE) --profile dev logs -f server-dev

.PHONY: docker-prod
docker-prod: ## Build production Docker image with multi-stage build
	@echo "$(GREEN)Building production Docker image...$(NC)"
	docker build -t $(BINARY_NAME):prod -f Dockerfile .
	@echo "$(GREEN)Production image built!$(NC)"
	@echo "Image size: $$(docker images $(BINARY_NAME):prod --format 'table {{.Size}}')"

.PHONY: docker-prod-run
docker-prod-run: ## Run production Docker container
	@echo "$(GREEN)Running production container...$(NC)"
	docker run -d \
		--name $(BINARY_NAME)-prod \
		-p 50051:50051 \
		--env-file .env \
		$(BINARY_NAME):prod
	@echo "$(GREEN)Production container running on port 50051!$(NC)"

.PHONY: docker-prod-stop
docker-prod-stop: ## Stop production Docker container
	@echo "$(YELLOW)Stopping production container...$(NC)"
	docker stop $(BINARY_NAME)-prod && docker rm $(BINARY_NAME)-prod
	@echo "$(GREEN)Production container stopped!$(NC)"

.PHONY: docker-clean
docker-clean: ## Remove all Docker resources for this project
	@echo "$(RED)Removing all Docker resources...$(NC)"
	$(DOCKER_COMPOSE) down -v --remove-orphans
	docker rmi -f $(BINARY_NAME):latest $(BINARY_NAME):prod 2>/dev/null || true
	@echo "$(GREEN)Docker cleanup complete!$(NC)"

.PHONY: db-start
db-start: ## Start PostgreSQL database container
	@echo "$(GREEN)Starting PostgreSQL container...$(NC)"
	docker run -d \
		--name sample-grpc-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=sample_grpc_server \
		-p 5432:5432 \
		postgres:15-alpine
	@echo "$(GREEN)PostgreSQL started!$(NC)"

.PHONY: db-stop
db-stop: ## Stop PostgreSQL database container
	@echo "$(YELLOW)Stopping PostgreSQL container...$(NC)"
	docker stop sample-grpc-postgres
	docker rm sample-grpc-postgres
	@echo "$(GREEN)PostgreSQL stopped!$(NC)"

.PHONY: db-shell
db-shell: ## Connect to PostgreSQL shell
	@echo "$(GREEN)Connecting to PostgreSQL...$(NC)"
	docker exec -it sample-grpc-postgres psql -U postgres -d sample_grpc_server

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)Tools installed!$(NC)"

.PHONY: check
check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "$(GREEN)All checks passed!$(NC)"

.PHONY: dev
dev: deps fmt ## Prepare for development
	@echo "$(GREEN)Development environment ready!$(NC)"

.PHONY: all
all: clean deps fmt vet lint test build ## Run all targets

.DEFAULT_GOAL := help