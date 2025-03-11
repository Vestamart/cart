BINARY_NAME=cart-service

build:
	go build -o $(BINARY_NAME) ./cmd


run:
	./$(BINARY_NAME)


run-all: build run


check-coverage:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | findstr total
	del coverage.out


test-bench:
	go test -bench=. ./internal/repository


cognitive-load:
	gocognit -top 10 -ignore "_mock|_test" .\internal


cyclomatic-load:
	gocyclo -top 10 -ignore "_mock|_test" .\internal