BINARY_NAME=cart-service

build:
	go build -o $(BINARY_NAME) ./cmd

run:
	./$(BINARY_NAME)

run-all: build run

test-coverage:
	go test -cover ./...
