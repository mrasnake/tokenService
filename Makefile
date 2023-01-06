build: ## Build server executable.
	go build -o ./cmd/api ./cmd/api

run: build ## Build and run server executable
	./cmd/api/api

help: build ## Run CLI help flag
	./cmd/api/api -h

test: ## Run the test suite.
	go test ./...