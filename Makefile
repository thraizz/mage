# MAGE Makefile
# 
# Main targets for building, testing, and deploying the MAGE Go server and web client.

-include .env

.PHONY: help
help:
	@echo "MAGE Build & Deployment Commands"
	@echo ""
	@echo "Development:"
	@echo "  make test              - Run all Go tests"
	@echo "  make test-go           - Run Go tests with coverage"
	@echo "  make build             - Build Go server binary"
	@echo "  make run               - Run Go server locally"
	@echo "  make proto             - Generate protobuf code"
	@echo ""
	@echo "Database:"
	@echo "  make db-up             - Start PostgreSQL with Docker"
	@echo "  make db-down           - Stop PostgreSQL"
	@echo "  make db-migrate        - Run database migrations"
	@echo "  make db-cards          - Download and import Scryfall cards"
	@echo ""
	@echo "Production:"
	@echo "  make deploy            - Deploy to production server"
	@echo "  make deploy-cards      - Update cards on production"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make clean-db          - Clean database (remove all data)"

# ============================================================================
# Development Targets
# ============================================================================

.PHONY: test
test:
	cd mage-server-go && go test ./...

.PHONY: test-go
test-go:
	cd mage-server-go && go test -v -race -coverprofile=coverage.out ./...
	cd mage-server-go && go tool cover -html=coverage.out -o coverage.html

.PHONY: build
build:
	cd mage-server-go && go build -o bin/mage-server ./cmd/web-demo/main.go

.PHONY: run
run:
	cd mage-server-go && go run ./cmd/web-demo/main.go

.PHONY: proto
proto:
	cd mage-server-go && ./scripts/generate_proto.sh

.PHONY: fmt
fmt:
	cd mage-server-go && go fmt ./...

.PHONY: lint
lint:
	cd mage-server-go && golangci-lint run ./...

# ============================================================================
# Database Targets
# ============================================================================

.PHONY: db-up
db-up:
	docker compose up -d postgres

.PHONY: db-down
db-down:
	docker compose down

.PHONY: db-migrate
db-migrate:
	cd mage-server-go && ./scripts/run_postgres_migrations.sh

.PHONY: db-cards
db-cards:
	cd mage-server-go && ./scripts/migrate_to_scryfall.sh

.PHONY: db-cards-download
db-cards-download:
	cd mage-server-go && ./scripts/download_scryfall_bulk.sh

# ============================================================================
# Production Deployment Targets
# ============================================================================

.PHONY: deploy
deploy:
	./deploy.sh

.PHONY: deploy-cards
deploy-cards:
	./update-cards-prod.sh --download

.PHONY: deploy-cards-test
deploy-cards-test:
	./update-cards-prod.sh --download --dry-run

# ============================================================================
# Cleanup Targets
# ============================================================================

.PHONY: clean
clean:
	cd mage-server-go && rm -rf bin/ tmp/ *.log *.out coverage.html
	cd mage-client-web && rm -rf .svelte-kit/ build/ node_modules/.cache/

.PHONY: clean-db
clean-db:
	docker compose down -v
	rm -f clean_dbs.sh || true

.PHONY: clean-all
clean-all: clean clean-db

# ============================================================================
# Docker Targets
# ============================================================================

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-logs
docker-logs:
	docker compose logs -f

.PHONY: docker-restart
docker-restart:
	docker compose restart

# ============================================================================
# Client Targets
# ============================================================================

.PHONY: client-install
client-install:
	cd mage-client-web && npm install

.PHONY: client-dev
client-dev:
	cd mage-client-web && npm run dev

.PHONY: client-build
client-build:
	cd mage-client-web && npm run build

.PHONY: client-preview
client-preview:
	cd mage-client-web && npm run preview

# ============================================================================
# All-in-one Targets
# ============================================================================

.PHONY: setup
setup: db-up db-migrate db-cards client-install
	@echo ""
	@echo "✓ Setup complete!"
	@echo ""
	@echo "To start development:"
	@echo "  1. Start backend:  make run"
	@echo "  2. Start frontend: make client-dev"

.PHONY: dev
dev:
	@echo "Starting development environment..."
	@echo "Backend: http://localhost:17171"
	@echo "Frontend: http://localhost:5173"
	@make -j2 run client-dev
