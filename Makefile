# Configuration
BIN_NAME      ?= amareh-bot
DOCKER_REG    ?= ghcr.io/sudosz/amareh
DOCKER_TAG    ?= latest
GO_PACKAGE    ?= ./cmd/bot
GO_TEST_FLAGS ?= -race -v -covermode=atomic -timeout=5m
GO_LDFLAGS    ?= -s -w -X main.version=$(shell git describe --tags --always)
GOLANGCI_LINT ?= golangci-lint

# Use := for immediate assignment where possible
BUILD_FLAGS := CGO_ENABLED=0 GO111MODULE=on
PLATFORMS := linux/amd64 linux/arm64
MAKEFLAGS += --jobs=$(shell nproc)

.PHONY: all build clean test lint format docker-build docker-push run setup help

# Default target with optimized dependency order
all: deps lint test build

## Build binary with optimizations and UPX compression
build:
	@printf "Building %s...\n" $(BIN_NAME)
	@mkdir -p bin/
	@$(BUILD_FLAGS) go build -trimpath -ldflags "$(GO_LDFLAGS)" -o bin/$(BIN_NAME) $(GO_PACKAGE)
	@which upx >/dev/null && upx --best --lzma bin/$(BIN_NAME) || true

## Build for multiple architectures efficiently
build-multi:
	@echo "Building for multiple platforms..."
	@mkdir -p bin/
	@$(foreach platform,$(PLATFORMS),\
		GOOS=$(firstword $(subst /, ,$(platform))) \
		GOARCH=$(lastword $(subst /, ,$(platform))) \
		$(BUILD_FLAGS) go build -trimpath -ldflags "$(GO_LDFLAGS)" \
			-o bin/$(BIN_NAME)-$(firstword $(subst /, ,$(platform)))-$(lastword $(subst /, ,$(platform))) $(GO_PACKAGE) & \
	)
	@wait

## Install and verify dependencies with caching
deps:
	@echo "Installing dependencies..."
	@go mod verify
	@go mod tidy -v
	@go mod download

## Run tests with optimized parallelization
test: deps
	@echo "Running tests..."
	@go test -count=1 $(GO_TEST_FLAGS) -parallel=$(shell nproc) -failfast ./...

## Run security checks with enhanced reporting
security:
	@echo "Checking for vulnerabilities..."
	@govulncheck -show verbose -tags=netgo,osusergo ./...

## Lint code with optimized settings
lint:
	@echo "Linting..."
	@$(GOLANGCI_LINT) run --fix --timeout 5m --max-same-issues 0 --max-issues-per-linter 0

## Format code with parallel execution
format:
	@echo "Formatting code..."
	@find . -name '*.go' -type f -print0 | xargs -0 -P $(shell nproc) gofmt -s -w
	@find . -name '*.go' -type f -print0 | xargs -0 -P $(shell nproc) goimports -w

## Build optimized Docker image with BuildKit
docker-build:
	@echo "Building Docker image..."
	@DOCKER_BUILDKIT=1 docker build --pull --no-cache --compress -t $(DOCKER_REG)/$(BIN_NAME):$(DOCKER_TAG) .

## Push Docker image with exponential backoff retry
docker-push:
	@echo "Pushing Docker image..."
	@for i in 1 2 4 8; do \
		docker push $(DOCKER_REG)/$(BIN_NAME):$(DOCKER_TAG) && break || \
		echo "Retrying in $$i seconds..." && sleep $$i; \
	done

## Run the application with optimized settings
run:
	@$(BUILD_FLAGS) go run -trimpath -gcflags='-N -l' $(GO_PACKAGE)

## Clean build artifacts and caches thoroughly
clean:
	@echo "Cleaning..."
	@rm -rf bin/ coverage.out
	@go clean -i -cache -testcache -modcache -fuzzcache

## Setup development environment with version checks
setup: deps
	@echo "Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go version

## Setup git hooks with validation
hooks:
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@install -m 755 scripts/pre-commit .git/hooks/pre-commit
	@test -x .git/hooks/pre-commit || chmod +x .git/hooks/pre-commit

## Show help with enhanced formatting
help:
	@echo "Available targets:"
	@awk -F':.*?## ' '/^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

## CI/CD setup with enhanced testing tools
setup-ci: deps
	@echo "Installing CI dependencies..."
	@go install gotest.tools/gotestsum@latest
	@go install github.com/axw/gocov/gocov@latest

## Generate code with validation
generate:
	@echo "Generating code..."
	@go generate ./...
	@go mod verify
	@go mod tidy