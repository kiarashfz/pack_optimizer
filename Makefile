.PHONY: dev-setup lint up down restart logs test

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

# Install necessary installation tools if they are not already installed
dev-setup: $(GOLANGCI_LINT) $(GOIMPORTS)

# The lint target now has both tools as prerequisites
lint: dev-setup
	@echo "Running goimports..."
	@find . -type f -name "*.go" | xargs -r $(GOIMPORTS) -w
	@echo "Running golangci-lint..."
	$(GOLANGCI_LINT) run ./...

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
