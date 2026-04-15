FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/project-api cmd/project-api/main.go

# ── Runtime ──────────────────────────────────────
FROM alpine:3.19

RUN adduser -D -s /bin/sh appuser
WORKDIR /app

COPY --from=builder /bin/project-api .

USER appuser
EXPOSE 8080

CMD ["./project-api"]
