# -------------------------------
# Configuration
# -------------------------------

APP_NAME=snippetbox
CMD_DIR=cmd/web
BIN_DIR=bin

# Go commands
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOVET=$(GO) vet
GOFMT=gofmt

# Database (used in CI / local if needed)
DB_NAME=snippetbox
DB_USER=root
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306

# -------------------------------
# Default target
# -------------------------------

.PHONY: all
all: test build

# -------------------------------
# Formatting
# -------------------------------

.PHONY: fmt
fmt:
	@echo ">> Formatting Go files"
	$(GOFMT) -w .

# -------------------------------
# Lint / Vet
# -------------------------------

.PHONY: vet
vet:
	@echo ">> Running go vet"
	$(GOVET) ./...

# -------------------------------
# Unit Tests (NO DB)
# -------------------------------

.PHONY: test
test:
	@echo ">> Running unit tests"
	$(GOTEST) ./...

# -------------------------------
# Integration Tests (DB)
# -------------------------------

.PHONY: test-integration
test-integration:
	@echo ">> Running integration tests"
	$(GOTEST) -tags=integration ./...
	
# -------------------------------
# Tests with coverage (CI-friendly)
# -------------------------------

.PHONY: test-cover
test-cover:
	@echo ">> Running tests with coverage"
	$(GOTEST) -cover ./...

# -------------------------------
# Build
# -------------------------------

.PHONY: build
build:
	@echo ">> Building application"
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)

# -------------------------------
# Clean
# -------------------------------

.PHONY: clean
clean:
	@echo ">> Cleaning build artifacts"
	rm -rf $(BIN_DIR)

# -------------------------------
# CI (used by GitHub Actions)
# -------------------------------

.PHONY: ci
ci: fmt vet test build
