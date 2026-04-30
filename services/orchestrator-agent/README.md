# Daedalus — Orchestrator Agent

Multi-agent orchestration backbone (Sprint 3 — PB-025).

The Orchestrator decomposes a high-level user goal (e.g. "Build a 200-unit/day
biscuit factory in Douala") into a directed graph of sub-tasks (procurement search,
layout generation, 3D render, cost analysis, …) and dispatches each one to the
correct specialised agent via **Redis Streams**.

## Architecture (Hexagonal / Kliops pattern)

```
cmd/orchestrator-api/        ← composition root + HTTP server
internal/
  core/
    domain/      ← Goal, SubTask, decomposition rules
    ports/       ← repositories + TaskPublisher interfaces
    services/    ← OrchestratorService (use-cases)
  adapters/
    handlers/    ← REST endpoints (Go 1.22 ServeMux patterns)
    repositories/← Postgres adapters (pgx)
    publishers/  ← Redis Streams adapter (XADD)
deployments/
  migrations/    ← SQL migrations
  docker-compose.yml
```

## Endpoints

| Method | Path                                  | Description                           |
|--------|---------------------------------------|---------------------------------------|
| POST   | `/api/orchestrator/goals`             | Submit a goal → returns task graph    |
| GET    | `/api/orchestrator/goals`             | List goals (filter `?user_id=`)       |
| GET    | `/api/orchestrator/goals/{id}`        | Get a goal + sub-tasks                |
| PATCH  | `/api/orchestrator/tasks/{id}`        | Update sub-task status                |
| GET    | `/health`                             | Health probe                          |
| GET    | `/metrics`                            | Prometheus metrics                    |

## Acceptance criteria — PB-025

> *Task graph visible in agent log; each subtask has unique ID and status.*

- ✅ Each sub-task has a UUID and a strict lifecycle (`pending` → `dispatched`
  → `in_progress` → `completed` / `failed`).
- ✅ Decomposition rules are encoded in `OrchestratorService.decompose`.
- ✅ The full task graph is logged on goal creation (`logTaskGraph`).
- ✅ A `TaskPublisher` port (`internal/core/ports/ports.go`) abstracts the
  dispatch transport. The default adapter is **InMemory** (great for dev
  and tests). A Redis Streams adapter (one stream per task type, e.g.
  `agent.procurement_search`) plugs in the same port — wire it in
  `cmd/orchestrator-api/main.go` once `github.com/redis/go-redis/v9` is
  available in your module cache.

## Run locally

```sh
make up      # docker-compose up postgres + redis
make migrate # run SQL migrations
make run     # go run ./cmd/orchestrator-api
```

Service listens on `:8082`.
