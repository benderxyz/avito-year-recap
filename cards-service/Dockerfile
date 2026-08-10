FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go build -o app ./cmd/main.go


FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN useradd --system --uid 10001 --no-create-home appuser

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/seeds ./seeds

USER appuser

EXPOSE 8081

CMD ["./app"]