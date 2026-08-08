FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8081

CMD ["./app"]