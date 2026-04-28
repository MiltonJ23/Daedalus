# Rapport — project-service

> SRS v2.0 §3.2 — CRUD projets factory (Sprint 1, **terminé**).

## Mission
Gestion du cycle de vie complet d'un projet d'usine : création, configuration
floor plan, archivage / restauration, suppression confirmée, dashboard,
auto-sauvegarde 60 s.

## Endpoints
| Méthode | Chemin                                | Description                        |
|---------|---------------------------------------|------------------------------------|
| POST    | `/api/projects`                       | Crée un projet (PB-014)            |
| GET     | `/api/projects?status=…`              | Liste filtrée (PB-016)             |
| GET     | `/api/projects/:id`                   | Détail                             |
| PUT     | `/api/projects/:id`                   | Mise à jour complète (PB-015)      |
| PATCH   | `/api/projects/:id/autosave`          | Patch incrémental (PB-017)         |
| PATCH   | `/api/projects/:id/archive`           | Archive / restore (PB-018)         |
| DELETE  | `/api/projects/:id?confirm=true`      | Suppression définitive             |
| GET     | `/health`                             | Liveness/readiness                 |
| GET     | `/metrics`                            | Prometheus                         |

## Architecture (Hexagonal)
```
cmd/project-api          # Entry point, DI, graceful shutdown
internal/
  core/
    domain/              # Entités, errors typées
    ports/               # Interfaces repository
    services/            # Use cases + validation
  adapters/
    handlers/            # HTTP (stdlib ServeMux), middleware, /metrics
    repositories/        # PostgreSQL (pgx/v5)
deployments/migrations/  # SQL versionné
```

## Tests
- 27 tests unitaires (services + handlers) — `go test -v ./...`
- E2E : `e2e-tests/` (Playwright + Jest).

## Métriques exposées (`/metrics`)
Cf. `services/project-service/internal/adapters/handlers/metrics.go` :
- `http_requests_total{method,path,status}`
- `http_request_duration_seconds{method,path}` (histogram)
- `http_requests_in_flight`
- `http_errors_total{method,path,status}`
- `project_operations_total{operation,result}`
