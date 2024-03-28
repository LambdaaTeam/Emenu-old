GO = go

build: ## Build the project
	$(GO) build -o bin/ ./...

clean: ## Clean the project
	rm -rf bin/

run-api: ## Run the API
	$(GO) run cmd/api/main.go

run-shortener: ## Run the Shortener
	$(GO) run cmd/shortener/main.go

run-ws: ## Run the WS Server
	$(GO) run cmd/ws/main.go

## Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[34m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)