.PHONY: help docs test lint tidy install

help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

test:
	@go test -v ./...

lint:
	@golangci-lint run

docs: ## Generate swagger documentation
	@swag init -g api.go

tidy:
	@go mod tidy

install: ## Install dependencies
	@GO111MODULE=off go get -u github.com/swaggo/swag/cmd/swag
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0
