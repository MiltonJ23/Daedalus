# Getting started

## Prérequis

* Docker 24+ et Docker Compose v2
* Go 1.23+ (dev local du `project-service`)
* Node.js 20+ (frontend)
* Python 3.11+ (services FastAPI/ADK)

## Lancement

```bash
cp .env.example .env
docker compose up --build
```

## Variables d'environnement

Toutes définies dans `.env` (cf. `.env.example`). Aucune duplication par
service : `docker-compose.yml` injecte le fichier via `env_file: .env`.

## Endpoints utiles

| URL                              | Description                  |
| -------------------------------- | ---------------------------- |
| http://localhost:3000            | Frontend Next.js             |
| http://localhost:8080/health     | Health check project-service |
| http://localhost:8080/metrics    | Métriques Prometheus         |
| http://localhost:9090            | Prometheus UI                |

## Tests

```bash
# Go
cd services/project-service && go test ./...

# Frontend
cd frontend && npm install && npm run lint && npm run build
```
