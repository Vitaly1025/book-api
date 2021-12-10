BINARY_NAME=book-api
TARGET_OS=windows
OUTPUT_DIR=${GOPATH}/bin

# Builds book api
build:
	GOARCH=amd64 GOOS=${TARGET_OS} go build -o '${OUTPUT_DIR}/${BINARY_NAME}-${TARGET_OS}' bookstore/cmd/bookstore-server

# Runs the database container
run_db:
	docker run --name=book-db -e POSTGRES_PASSWORD='1234qwerty' -p 5436:5432 -d --rm postgres
# Runs book api
run:
	GOARCH=amd64 GOOS=${TARGET_OS} go run book-api/cmd/main.go

# Runs all
make run_all:
	docker-compose up

# Checks sources
lint:
	golangci-lint run -v

# Test book api
.PHONY: test
test:
	go test .\pkg\repository && go test .\pkg\service
# Migrates the database UP
migration_up:
	migrate -path ./schema -database 'postgres://postgres:1234qwerty@db:5432/postgres?sslmode=disable' up
# Migrates the database DOWN
migration_down:
	migrate -path ./schema -database 'postgres://postgres:1234qwerty@db:5432/postgres?sslmode=disable' down

# Downloads modules
dep:
	go mod download
	
# Downloads additional tools
prepare:
	go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate/