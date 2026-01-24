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
all: ci

# -------------------------------
# Formatting (Check only - CI safe)
# -------------------------------

.PHONY: fmt
fmt:
	@echo ">> Checking Go formatting"
	test -z "$$($(GOFMT) -l .)"

# Optional: developer-only formatter
.PHONY: fmt-fix
fmt-fix:
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
# Run unit tests only (NO DB)
# -------------------------------

.PHONY: test-unit
test-unit:
	@echo ">> Running unit tests only"
	$(GOTEST) ./... -short

# -------------------------------
# Run all tests (unit + integration(DB) tests)
# (NOT USED in Stage 2 CI)
# -------------------------------

.PHONY: test-integration
test-integration:
	@echo ">> Running unit + integration tests"
	$(GOTEST) ./...


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
# CI (Stage 2)
# -------------------------------

.PHONY: ci
ci: fmt vet test-unit build