# procurement-agent

Daedalus Procurement Agent — Go microservice that searches global supplier
catalogs, ranks equipment using a composite score
(price 40% / spec match 35% / supplier reliability 15% / lead time 10%),
caches results for 24h, and lets users approve / reject recommendations.

Implements requirements **FR-PROC-01 → FR-PROC-07** from the SRS.

## Architecture

Hexagonal (ports & adapters) — same layout as `project-service`:

```
cmd/procurement-api/        # entry point (HTTP server)
internal/
  core/
    domain/                 # pure entities & errors
    ports/                  # repository & supplier interfaces
    services/               # use-case orchestration + ranking
  adapters/
    handlers/               # HTTP handlers + middleware + Prometheus
    repositories/           # PostgreSQL (pgx) persistence
    suppliers/              # external supplier catalog adapters
deployments/                # docker-compose + SQL migrations
```

## Endpoints

| Method | Path                                              | Description                                  |
|--------|---------------------------------------------------|----------------------------------------------|
| POST   | `/api/procurement/searches`                       | Submit an equipment search request           |
| GET    | `/api/procurement/searches/{id}`                  | Fetch a search and its ranked results        |
| GET    | `/api/procurement/searches`                       | List searches (filter by `?project_id=`)     |
| PATCH  | `/api/procurement/results/{id}/decision`          | Approve or reject a recommendation           |
| GET    | `/health`                                         | Liveness / readiness probe                   |
| GET    | `/metrics`                                        | Prometheus metrics                           |

## Run

```bash
make docker-up        # start PostgreSQL
make build && make run
```

Default port: `8081` (override with `PORT`). DB DSN via `DB_DSN`.

## Tests

```bash
make test
```
