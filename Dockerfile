FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cart-service ./cmd/server


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/cart-service .
EXPOSE 8082
CMD ["./cart-service"]