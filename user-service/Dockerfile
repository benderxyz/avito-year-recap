FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

ENV MIGRATIONS_DIR=/app/migrations
ENV SERVER_PORT=8082

EXPOSE 8082

CMD ["./app"]
