.PHONY: setup-docker
setup-docker: ## Setup docker resources for testing
	docker-compose up -d

.PHONY: clean-docker
clean-docker: ## Clean docker resources for testing
	docker-compose down -v

.PHONY: test
test: ## Run unit tests
	go test -v ./...

.PHONY: test-integration
test-integration: ## Rull all tests including integration tests
	go test -v -p 1 -tags integration ./...

.PHONY: test-generate
test-command: ## Run the command to test generating codes.
	./scripts/test-command.sh

.PHONY: test-ci
test-ci: ## Run tests for CI
	go test -v -p 1 -race -coverprofile coverage.coverprofile -tags integration ./...

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
