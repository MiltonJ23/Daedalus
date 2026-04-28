# Daedalus

> SaaS de conception architecturale assistée par IA — génération de plans 2D/3D,
> sourcing automatisé et estimation budgétaire (SRS v2.0).

[![License](https://img.shields.io/badge/license-Proprietary-blue.svg)](LICENSE.md)

---

## 🚀 Démarrage rapide

```bash
git clone https://github.com/MiltonJ23/Daedalus.git
cd Daedalus
cp .env.example .env          # ajuster les secrets
docker compose up --build
```

| Service        | URL                          |
| -------------- | ---------------------------- |
| Frontend       | http://localhost:3000        |
| Project API    | http://localhost:8080        |
| Metrics        | http://localhost:8080/metrics |
| Prometheus     | http://localhost:9090        |

Plus de détails : [`wiki/getting-started.md`](wiki/getting-started.md).

---

## 🏗️ Architecture

Microservices polyglottes (Go + Python) orchestrés derrière Traefik avec
RabbitMQ comme bus d'événements et Redis pour le cache des quotas.

```
frontend (Next.js)
    │
    ▼
Traefik ── ForwardAuth ── auth-service (Go)
    │
    ├── project-service     (Go)
    ├── payment-service     (Python / FastAPI)
    ├── notification-service(Python)
    ├── orchestrator-agent  (Python / ADK)
    ├── procurement-agent   (Python / ADK + Playwright)
    ├── layout-engine       (Python / ADK)
    ├── 3d-asset-service    (Python / Meshy.ai)
    └── analytics-service   (Python / FastAPI)
```

Détails : [`architecture/project-management-service/README.md`](architecture/project-management-service/README.md)
et [`wiki/architecture.md`](wiki/architecture.md).

---

## 📁 Organisation du dépôt

| Dossier              | Rôle                                               |
| -------------------- | -------------------------------------------------- |
| `services/`          | Code source des microservices                      |
| `frontend/`          | Application Next.js 14 (R3F pour la 3D)            |
| `shared/`            | Bibliothèques partagées (logging, etc.)            |
| `API/`               | Spécifications OpenAPI 3.0 par service             |
| `ci-cd/`             | Pipelines Jenkins par service                      |
| `rapport/`           | Documentation fonctionnelle par service            |
| `architecture/`      | Diagrammes C4, séquence, ADRs                      |
| `infrastructure/`    | Terraform, Ansible, Kubernetes, Prometheus/Grafana |
| `wiki/`              | Documentation transverse                           |
| `e2e-tests/`         | Tests bout-en-bout (Playwright)                    |
| `Daedalus_SRS_v2.0.pdf` | Spécifications fonctionnelles                  |

---

## 🧪 Qualité & Observabilité

* **Tests Go** : `cd services/project-service && go test ./...`
* **Linting** : `golangci-lint run` (Go) · `ruff` (Python) · `npm run lint` (frontend)
* **Métriques** : tous les services exposent `/metrics` (Prometheus)
* **Tracing** : OpenTelemetry → Jaeger
* **CI/CD** : Jenkins (cf. [`ci-cd/README.md`](ci-cd/README.md))

---

## 🤝 Contribuer

Voir [`CONTRIBUTING.md`](CONTRIBUTING.md) et [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).

## 📄 Licence

Proprietary — voir [`LICENSE.md`](LICENSE.md).
