.PHONY: lint up down restart logs test

# Locate the golangci-lint executable
GOLANGCI_LINT := $(shell which golangci-lint 2>/dev/null || echo $(shell go env GOPATH)/bin/golangci-lint)
# Locate the goimports executable
GOIMPORTS := $(shell which goimports 2>/dev/null || echo $(shell go env GOPATH)/bin/goimports)

# Rule to install golangci-lint if it's not found.
$(GOLANGCI_LINT):
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Rule to install goimports if it's not found
$(GOIMPORTS):
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest

# The lint target now has both tools as prerequisites
lint: $(GOLANGCI_LINT) $(GOIMPORTS)
	@echo "Running goimports..."
	@find . -type f -name "*.go" | xargs -r $(GOIMPORTS) -w
	@echo "Running golangci-lint..."
	$(GOLANGCI_LINT) run ./...


VENV_DIR := .venv
dev-setup:
	@echo "Checking for Python 3..."
	@command -v python3 >/dev/null 2>&1 || { echo >&2 "Python 3 is not installed. Please install Python 3."; exit 1; }
	@echo "Setting up Python virtual environment..."
	@test -d $(VENV_DIR) || python3 -m venv $(VENV_DIR)
	@echo "Upgrading pip in virtual environment..."
	@$(VENV_DIR)/bin/python -m pip install --upgrade pip
	@echo "Installing pre-commit..."
	@$(VENV_DIR)/bin/pip install pre-commit
	@echo "Installing pre-commit hooks..."
	@$(VENV_DIR)/bin/pre-commit install
	@echo "✅ Development environment setup complete."


# Start all services defined in docker-compose
up:
	docker compose up -d --build

# Stop and remove containers, networks, volumes, and images created by up
down:
	docker compose down

# Restart services
restart: down up

# Follow logs for all services
logs:
	docker compose logs -f

# Run all Go tests in the 'test' directory and its subdirectories
test:
	go test -v -race ./test/...