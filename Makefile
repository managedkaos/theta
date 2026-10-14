BINARY_NAME := theta
BUILD_DIR := bin
MODULE := github.com/managedkaos/theta
VERSION ?= $(shell git describe --tags --dirty=-local --match 'v[0-9]*' 2>/dev/null || echo dev)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILD_ENV ?= local
PREFIX ?= $(HOME)/.local
LDFLAGS := -ldflags "-s -w -X '$(MODULE)/cmd.Version=$(VERSION)' -X '$(MODULE)/cmd.GitBranch=$(GIT_BRANCH)' -X '$(MODULE)/cmd.BuildDate=$(BUILD_DATE)' -X '$(MODULE)/cmd.BuildEnvironment=$(BUILD_ENV)'"

.PHONY: help build build-all test lint verify run man snapshot install
help: ## Display available targets
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*## / {printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the native binary in bin/theta
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

build-all: ## Build Linux, macOS, and Windows binaries
	@mkdir -p $(BUILD_DIR)
	@set -e; for platform in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do \
		os_name=$${platform%/*}; arch_name=$${platform#*/}; ext=""; \
		if [ "$$os_name" = windows ]; then ext=.exe; fi; \
		CGO_ENABLED=0 GOOS=$$os_name GOARCH=$$arch_name go build -trimpath $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$$os_name-$$arch_name$$ext .; \
	done

test: ## Run tests with race detection
	go test -race ./...

lint: ## Check formatting and run go vet
	@test -z "$$(gofmt -l main.go cmd internal tools)" || { gofmt -l main.go cmd internal tools; exit 1; }
	go vet ./...

verify: lint test build ## Check, test, and build theta

run: ## Run theta with ARGS
	go run . $(ARGS)

man: ## Generate man pages
	go run ./tools/manpages

snapshot: ## Build release archives without publishing
	goreleaser release --snapshot --clean

install: build man ## Install binary and man page under PREFIX
	install -d "$(DESTDIR)$(PREFIX)/bin" "$(DESTDIR)$(PREFIX)/share/man/man1"
	install -m 0755 $(BUILD_DIR)/$(BINARY_NAME) "$(DESTDIR)$(PREFIX)/bin/$(BINARY_NAME)"
	install -m 0644 man/theta.1 "$(DESTDIR)$(PREFIX)/share/man/man1/theta.1"
