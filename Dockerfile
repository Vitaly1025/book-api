#download go and install dependencies
FROM golang:alpine AS base
WORKDIR /src
COPY go.* .
RUN go mod download
COPY . .

#make build with BuildKit
FROM base AS build
ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \ 
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /book-api ./cmd/main.go

#launch tests
FROM base AS unit-test
RUN --mount=type=cache,target=/root/.cache/go-build \
go test .\pkg\repository && go test .\pkg\service

#set scratch system
FROM scratch 
EXPOSE 4000
COPY --from=build /book-api /
