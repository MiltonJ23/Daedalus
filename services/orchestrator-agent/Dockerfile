# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/orchestrator-api ./cmd/orchestrator-api

FROM alpine:3.20
RUN adduser -D -u 10001 daedalus
COPY --from=builder /out/orchestrator-api /usr/local/bin/orchestrator-api
USER daedalus
EXPOSE 8082
ENTRYPOINT ["/usr/local/bin/orchestrator-api"]
