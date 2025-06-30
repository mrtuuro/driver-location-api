APP_NAME = driver-location-api

GO      = go
CMD_DIR = ./cmd
BIN     = bin/$(APP_NAME)
MAIN    = $(CMD_DIR)/main.go
ENV     = .env

SWAG_CLI = swag
SWAG_MAIN = cmd/main.go
SWAG_OUT = docs

.PHONY: all swagger build run test clean fmt lint tidy

all: build

swagger: 
	@$(SWAG_CLI) init -g $(SWAG_MAIN) --output $(SWAG_OUT)


build: swagger
	@mkdir -p bin
	@$(GO) build -o $(BIN) $(MAIN)

run: build
	@$(ENV_LOADER) $(BIN)

test:
	$(GO) test ./... -v -cover

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin

env:
	@echo "Printing .env config:"
	@cat $(ENV)

import:
	@go run ./tools/importer/main.go

help:
	@echo "Makefile Commands:"
	@echo "  build     - Build the Go binary"
	@echo "  run       - Build and run the binary"
	@echo "  test      - Run unit tests"
	@echo "  clean     - Remove build artifacts"
	@echo "  env       - Show current .env variables"
	@echo "  import    - Insert driver locations as batch"
