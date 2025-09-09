# Variables
GO_FILES := $(shell find . -name '*.go')
VERSION := $(shell git describe --tags --always --dirty)
APP_NAME := openapi-converter

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist

# Default target
all: build

# Build the Go application
build: $(GO_FILES)
	@echo "Building $(APP_NAME)..."
	@go build -o $(GOBIN)/$(APP_NAME) ./cmd

# Run the built application
run:
	@echo "Running $(APP_NAME)..."
	@$(GOBIN)/$(APP_NAME)

# Run the current code directly
dev:
	@echo "Running $(APP_NAME) in dev mode..."
	@go run ./cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	@cd test/go && go test -v ./...

# Build and install to system
install: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	@cp $(GOBIN)/$(APP_NAME) /usr/local/bin/

# Runs SBump, updates the patch version, commits, tags, and pushes
deploy:
	@echo "Bumping version and deploying"
	@./sbump.sh patch --push-version
	@echo -e "New version deploying /n"

.PHONY: build run dev test install deploy all