# 🏗️ Daedalus — Industrial Digital Twins Platform

> Sprint 2 : Project Management & CI/CD (34 SP)

## Architecture

Go Clean Architecture (Hexagonal), following Kliops reference project:

```
services/project-service/
├── cmd/project-api/     # Entry point, DI wiring, graceful shutdown
├── internal/
│   ├── core/
│   │   ├── domain/      # Entities, error types
│   │   ├── ports/       # Repository interfaces
│   │   └── services/    # Use cases, validation
│   └── adapters/
│       ├── handlers/    # HTTP handlers (stdlib ServeMux), middleware
│       └── repositories/# PostgreSQL (pgx/v5)
├── deployments/
│   ├── docker-compose.yml
│   └── migrations/
├── Makefile
├── Dockerfile
└── go.mod
```

## Quick Start

### Prerequisites
- Go 1.22+, Node.js 20+, Docker, Docker Compose

### Local Development (Docker Compose)

```bash
docker-compose up --build
```

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/api
- **Health check**: http://localhost:8080/health

### Manual Setup

```bash
# Backend
cd services/project-service
cp .env.example .env   # configure DB_DSN
make docker-up          # start PostgreSQL
psql $DB_DSN -f deployments/migrations/001_create_projects.sql
make run                # builds and runs on :8080

# Frontend
cd frontend
npm install
NEXT_PUBLIC_API_URL=http://localhost:8080/api npm run dev
```

### Run Tests

```bash
cd services/project-service
make test     # go test -v ./...
make vet      # go vet ./...
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/projects` | Create project |
| GET | `/api/projects` | List projects (`?status=active\|archived\|all`) |
| GET | `/api/projects/:id` | Get project |
| PUT | `/api/projects/:id` | Update project |
| PATCH | `/api/projects/:id/autosave` | Auto-save (incremental) |
| PATCH | `/api/projects/:id/archive` | Archive / restore (`?action=restore`) |
| DELETE | `/api/projects/:id?confirm=true` | Permanent delete |
| GET | `/health` | Health check |
| GET | `/metrics` | Prometheus metrics |

Full OpenAPI spec: `docs/api/swagger.yaml`

## Sprint 2 — User Stories

| ID | Story | SP | Status |
|----|-------|----|--------|
| PB-011 | Jenkins pipeline with GitHub integration | 8 | ✅ |
| PB-012 | Docker build & push to registry | 5 | ✅ |
| PB-013 | Auto-deploy to Kubernetes (K3s) | 5 | ✅ |
| PB-014 | Create factory project (CRUD) | 5 | ✅ |
| PB-015 | Factory floor dimensions | 3 | ✅ |
| PB-016 | Dashboard with status indicators | 3 | ✅ |
| PB-017 | Auto-save every 60 seconds | 3 | ✅ |
| PB-018 | Archive / delete projects | 2 | ✅ |

## Testing

```bash
# Backend tests (27 tests — service + handler)
cd services/project-service
go test -v ./...

# Backend lint
go vet ./...

# Frontend lint + build
cd frontend
npm run lint && npm run build

# E2E tests (requires running stack)
npx jest e2e-tests/
```

## CI/CD Pipeline (Jenkins)

Location: `ci-cd/jenkins/Jenkinsfile`

```
Push → Lint & Test → Docker Build → Push Registry → Deploy K8s
```

- Triggers on every GitHub push
- Parallel Docker builds (backend + frontend)
- Auto-deploy with `kubectl rollout status` verification

## Kubernetes

```bash
kubectl apply -f infrastructure/kubernetes/
kubectl rollout status deployment/backend -n daedalus
```

## Monitoring

- **Prometheus**: scrapes `/metrics` → `infrastructure/monitoring/prometheus.yml`
- **Grafana**: dashboard → `infrastructure/monitoring/grafana-dashboard.json`

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22+, stdlib http.ServeMux, pgx/v5 |
| Frontend | Next.js 14, React 18, TypeScript 5, Tailwind CSS |
| Database | PostgreSQL 16 |
| CI/CD | Jenkins, Docker, Kubernetes (K3s) |
| Monitoring | Prometheus, Grafana |
