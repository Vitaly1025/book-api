# download go and install dependencies
FROM golang:alpine AS base
WORKDIR /src
COPY go.* .
RUN go mod download
COPY . .

# make build with BuildKit
FROM base AS build
WORKDIR /src
ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/book-api/main.go
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o /book-api ./cmd/book-api/main.go 

# launch tests
FROM base AS unit-test
WORKDIR /src
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test ./internal/repository && go test ./internal/service 

# set scratch system
FROM alpine 
WORKDIR /app
COPY --from=build /book-api /app/book-api
CMD ["/app/book-api", "-config", "/app/config.yaml"]