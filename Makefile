BINARY_NAME=cart-service

build:
	go build -o $(BINARY_NAME) ./cmd

run:
	./$(BINARY_NAME)

run-all: build run

test-coverage:
	go test -cover ./...


check-coverage:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | findstr total
	del coverage.out