FROM golang:1.17.3-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install tool for migration
RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate/

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o book-api ./cmd/main.go

EXPOSE 4000

CMD [ "./book-api" ]

