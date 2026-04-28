# Architecture — Project Management Service

> SRS v2.0 §3.2 — Sprint 1 (terminé).

## Vue d'ensemble
Service Go autonome appliquant **Clean / Hexagonal Architecture** (référence
Kliops). Responsable du cycle de vie complet d'un projet d'usine
(création, configuration floor plan, archivage, suppression confirmée,
auto-sauvegarde 60 s, dashboard).

```
┌─────────────────────────────────────────────────────────┐
│                       cmd/project-api                    │
│        (DI wiring, graceful shutdown, /health, /metrics) │
└──────────────┬──────────────────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────────────────┐
│                internal/adapters/handlers                │
│   HTTP (stdlib ServeMux), CORS, logging, Prometheus mw   │
└──────────────┬──────────────────────────────────────────┘
               │ dépend de ports.ProjectRepository
               ▼
┌─────────────────────────────────────────────────────────┐
│                 internal/core/services                   │
│             use cases + validations métier               │
└──────────────┬──────────────────────────────────────────┘
               │ dépend de ports
               ▼
┌─────────────────────────────────────────────────────────┐
│           internal/adapters/repositories                 │
│              PostgreSQL (pgx/v5 pool)                    │
└──────────────┬──────────────────────────────────────────┘
               │
               ▼
                    ┌──────────────┐
                    │ PostgreSQL 16│
                    └──────────────┘
```

## Couches (Hexagonal)

| Couche                                | Rôle                                                         |
|---------------------------------------|--------------------------------------------------------------|
| `internal/core/domain/`               | Entités (`Project`), erreurs typées (`ErrNotFound`, etc.)    |
| `internal/core/ports/`                | Interfaces : `ProjectRepository`                             |
| `internal/core/services/`             | Use cases, règles métier, validation                         |
| `internal/adapters/handlers/`         | HTTP handlers, middleware (CORS, log, **metrics**)            |
| `internal/adapters/repositories/`     | Implémentation PostgreSQL des ports                          |

L'inversion de dépendances permet :
- Tests unitaires sans DB (mock du `ProjectRepository`).
- Substitution future (ex. cache Redis devant Postgres).

## Endpoints (SRS §3.2 + PB-014 → PB-018)
| Méthode | Chemin                                | Rôle                          |
|---------|---------------------------------------|-------------------------------|
| POST    | `/api/projects`                       | Créer (PB-014)                |
| GET     | `/api/projects?status=…`              | Lister (PB-016)               |
| GET     | `/api/projects/:id`                   | Détail                        |
| PUT     | `/api/projects/:id`                   | Mise à jour (PB-015)          |
| PATCH   | `/api/projects/:id/autosave`          | Auto-save 60 s (PB-017)       |
| PATCH   | `/api/projects/:id/archive`           | Archive/restore (PB-018)      |
| DELETE  | `/api/projects/:id?confirm=true`      | Suppression définitive        |
| GET     | `/health`                             | K8s liveness/readiness        |
| GET     | `/metrics`                            | Prometheus (NFR-O01)          |

## Modèle de données (extrait)
```sql
CREATE TABLE projects (
  id          UUID PRIMARY KEY,
  user_id     UUID NOT NULL,
  name        TEXT NOT NULL,
  industry    TEXT,
  floor       JSONB NOT NULL,        -- {l, w, h, units}
  status      TEXT NOT NULL DEFAULT 'active', -- active|archived
  budget_xaf  NUMERIC,
  created_at  TIMESTAMPTZ DEFAULT now(),
  updated_at  TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_projects_user_status ON projects(user_id, status);
```

## Intégrations
- **Auth** : reçoit headers `X-User-Id`, `X-User-Role`, `X-User-Plan` via
  Traefik Forward Auth (NFR-SEC04). Le **plan** sert à enforcer la limite
  projets actifs (Free 1 / Starter 5 / Business 20 / Enterprise ∞).
- **Notification** : publie `project.created`, `project.archived` sur
  RabbitMQ → mails transactionnels.
- **Layout / 3D / Analytics** : consomment `project_id` pour leurs jobs.

## Observabilité
Métriques Prometheus exposées sur `/metrics` (cf. `internal/adapters/handlers/metrics.go`) :
- `http_requests_total{method,path,status}`
- `http_request_duration_seconds{method,path}` (histogram)
- `http_requests_in_flight` (gauge)
- `http_errors_total{method,path,status}`
- `project_operations_total{operation,result}` (create/update/archive/delete)

Traces Jaeger via OpenTelemetry (planifié Sprint 6, NFR-O01).

## Diagrammes
Voir `diagrams/` :
- `c4-context.md`           — C4 niveau 1 (système & acteurs)
- `c4-container.md`         — C4 niveau 2 (containers internes)
- `sequence-create-project.md` — séquence "Create Project"
- `sequence-autosave.md`    — séquence auto-save 60 s
