GO=go

.DEFAULT_GOAL := help

build: ## Build the project
	$(GO) build -o bin/ ./...

test: ## Run the tests
	$(GO) test -v ./...

clean: ## Clean the project
	rm -rf bin/

install: ## Install project dependencies
	$(GO) mod download && go mod verify

run-api: ## Run the API
	$(GO) run ./cmd/api

run-shortener: ## Run the Shortener
	$(GO) run ./cmd/shortener

run-ws: ## Run the WS Server
	$(GO) run ./cmd/ws

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[34m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)