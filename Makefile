BINARY_NAME_CART=cart-service


build-cart:
	go build -o $(BINARY_NAME_CART) ./cmd/server

run-cart:
	./$(BINARY_NAME_CART)



run-all: build-cart run-cart

# Проверка покрытия
check-coverage:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | findstr total
	rm coverage.out

# Бенчмарки
test-bench:
	go test -bench=. ./internal/repository

# Когнитивная нагрузка
cognitive-load:
	gocognit -top 10 -ignore "_mock|_test" .\internal

# Цикломатическая нагрузка
cyclomatic-load:
	gocyclo -top 10 -ignore "_mock|_test" .\internal

# Контроль рейзов
race-check:
	go run -race ./cmd/server/main.go

