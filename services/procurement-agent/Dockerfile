FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/procurement-api cmd/procurement-api/main.go

# ── Runtime ──────────────────────────────────────
FROM alpine:3.19

RUN adduser -D -s /bin/sh appuser
WORKDIR /app

COPY --from=builder /bin/procurement-api .

USER appuser
EXPOSE 8081

CMD ["./procurement-api"]
