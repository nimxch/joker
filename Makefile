# Makefile for Go application

APP_NAME := joker
GO_FILES := $(shell find . -type f -name '*.go')
BUILD_DIR := .build
BIN_PATH := $(BUILD_DIR)/$(APP_NAME)

.PHONY: all build clean run test go-run

all: build

build: $(GO_FILES)
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BIN_PATH) .
dev:
	@go run .

run: build
	@echo "Running $(APP_NAME)..."
	@$(BIN_PATH)

go-run:
	@echo "Running $(APP_NAME) with 'go run'..."
	@go run $(GO_FILES)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)