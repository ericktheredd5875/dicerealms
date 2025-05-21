# Makefile for DiceRealms

APP_NAME=dicerealms

# Detect OS and set binary appropriately
ifeq ($(OS),Windows_NT)
	BINARY_EXT=.exe
else
	BINARY_EXT=
endif

BINARY=bin/$(APP_NAME)$(BINARY_EXT)

ifneq (,$(wildcard .env))
	-include .env
	export
endif

.PHONY: all build run test tidy clean dev

# Run go mod tidy
tidy:
	go mod tidy

# Build the Go binary
build:
	go build -o $(BINARY) ./cmd/server

# Run the server (Requires DB to be running)
run:
	go run ./cmd/server

# Run tests
test:
	go test ./...

test-reset:
	migrate -path ./db/migrations -database "$(TEST_DATABASE_URL)" -verbose force 3 || true
	migrate -path ./db/migrations -database "$(TEST_DATABASE_URL)" -verbose up
	go test ./...

# Clean the build binary
clean: 
	rm -rf bin tmp

# Run dev server with auto-reload
dev:
	air -c .air.toml

# migrate -path ./db/migrations -database "$(DATABASE_URL)" -verbose force 1 || true
migrate:
	migrate -path ./db/migrations -database "$(DATABASE_URL)" -verbose force 3 || true
	migrate -path ./db/migrations -database "$(DATABASE_URL)" -verbose up

# Run tests with coverage and save to coverage.out
coverage:
	go test ./... --coverprofile=coverage.out

# Show coverage % from coverage.out
coverage-show: coverage
	go tool cover -func=coverage.out

# View HTML coverage report
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

# Run on Render
render: 
	migrate
	./dicerealms