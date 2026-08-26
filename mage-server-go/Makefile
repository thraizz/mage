.PHONY: build test bot-sim proto run clean docker lint fmt help

# Variables
BINARY_NAME=mage-server
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
PROTO_DIR=./api/proto
PKG_PROTO_DIR=./pkg/proto

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

# Bot simulation (make bot-sim)
BOT_SIM_GAMES?=100
BOT_SIM_TIMEOUT?=30m

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: build ## Build and run the server
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME) -config config/config.yaml

run-log: build ## Build and run server with output logged to server.log
	@echo "Running $(BINARY_NAME) with logging to server.log..."
	$(BUILD_DIR)/$(BINARY_NAME) -config config/config.yaml 2>&1 | tee server.log

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Test coverage:"
	$(GOCMD) tool cover -func=coverage.out

bot-sim: ## Run N headless 4-bot games and report the completion rate (N=BOT_SIM_GAMES, default 100)
	@echo "Running $(BOT_SIM_GAMES) headless bot games..."
	BOT_SIM_GAMES=$(BOT_SIM_GAMES) $(GOTEST) -count=1 -timeout $(BOT_SIM_TIMEOUT) -v \
		-run TestBotSimCompletionRate ./internal/bot/

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GOTEST) -v -race -tags=integration ./test/integration/...

proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	@./scripts/generate_proto.sh
	@echo "Protobuf generation complete"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@echo "Clean complete"

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@goimports -w .

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t mage-server:latest .

docker-run: docker-build ## Run Docker container
	@echo "Running Docker container..."
	docker-compose up

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path migrations -database "postgresql://mage:mage@localhost:5432/mage?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	migrate -path migrations -database "postgresql://mage:mage@localhost:5432/mage?sslmode=disable" down

migrate-create: ## Create new migration (usage: make migrate-create NAME=create_users)
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir migrations -seq $(NAME)

wipe-sessions: ## Wipe all active game sessions from the database
	@echo "Wiping active game sessions..."
	PAGER=cat psql "postgresql://mage:mage@localhost:5432/mage?sslmode=disable" -c "TRUNCATE TABLE active_games;"
	@echo "Game sessions wiped"


transpile-cards: ## Transpile cards
	@echo "Transpiling cards..."
	python3 ./scripts/transpile_cards.py --batch
	go fmt ./internal/game/cards/generated/...
	@echo "Cards transpiled"

tools: ## Install development tools
	@echo "Installing development tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installed"

.DEFAULT_GOAL := help
