APP_NAME := book-api
BUILD_DIR := build
VERSION := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
MAIN := ./cmd/$(APP_NAME)/main.go
LDFLAGS := -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'"

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Builds book api
build:
	go build -o bin/$(APP_NAME) $(MAIN) t

# Runs the database container
docker-build:
	docker build -t $(APP_NAME) .

# Checks sources
lint:
	golangci-lint run -v

# Downloads modules
dep:
	go mod download

# Runs book api
run-locally:
	go run $(MAIN)

# Runs all
compose-up:
	docker compose up db book-api -d

compose-down:
	docker compose down

# Test book api
test:
	go test ./internal/repository && go test ./internal/service
	
# Migrates the database UP
migration_up:
	docker compose up migrator
	
# Downloads additional tools
prepare:
	go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate/