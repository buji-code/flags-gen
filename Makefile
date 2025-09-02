# Makefile for flags-gen

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Binary info
BINARY_NAME=flags-gen
BINARY_PATH=./bin/$(BINARY_NAME)
CMD_PATH=./cmd/$(BINARY_NAME)

# Build info
VERSION ?= dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

.PHONY: all build clean test coverage deps help install lint fmt vet

all: test build

## Build the binary
build:
	mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(CMD_PATH)

## Build for multiple platforms
build-all:
	mkdir -p dist
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)

## Install the binary to $GOPATH/bin
install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(CMD_PATH)

## Run tests
test:
	$(GOTEST) -v ./...

## Run tests with race detection
test-race:
	$(GOTEST) -race -v ./...

## Run tests with coverage
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem ./...

## Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html

## Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## Format code
fmt:
	$(GOFMT) -s -w .

## Run go vet
vet:
	$(GOCMD) vet ./...

## Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "Please install golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

## Run security scan (requires gosec)
security:
	@which gosec > /dev/null || (echo "Please install gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...

## Check for vulnerabilities (requires govulncheck)
vuln:
	@which govulncheck > /dev/null || (echo "Please install govulncheck: go install golang.org/x/vuln/cmd/govulncheck@latest" && exit 1)
	govulncheck ./...

## Run all quality checks
check: fmt vet lint test

## Generate example
example:
	$(BINARY_PATH) -i internal/testdata/example.go -o internal/testdata/example_flags.go

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-15s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Development targets

## Watch for changes and rebuild
dev:
	@which entr > /dev/null || (echo "Please install entr for file watching" && exit 1)
	find . -name '*.go' | entr -r make build

## Run the tool with example
run-example: build
	$(BINARY_PATH) -i internal/testdata/example.go

## Quick development test
quick: fmt vet build test