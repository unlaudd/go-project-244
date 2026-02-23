.PHONY: build lint test

build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run

test:
	go test -v -cover ./...
