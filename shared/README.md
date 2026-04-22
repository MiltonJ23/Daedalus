# Daedalus — Shared Module

Ce répertoire contient toutes les configurations, abstractions et utilitaires **partagés entre tous les microservices** du monorepo Daedalus. Il est structuré pour supporter à la fois les services **Go** et les services **Python (FastAPI / ADK)**.

---

## Structure du répertoire

```
shared/
├── go/                          # Packages Go partagés
│   ├── config/                  # Chargement centralisé de config (env vars, Viper)
│   ├── errors/                  # Erreurs métier personnalisées (codes, messages)
│   ├── middleware/              # Middleware HTTP (auth JWT, logging, rate-limit)
│   ├── messaging/               # Client RabbitMQ (publish/subscribe abstraits)
│   ├── cache/                   # Client Redis (get, set, del, pub/sub)
│   ├── storage/                 # Client MinIO (upload, download, presigned URLs)
│   ├── tracing/                 # Setup OpenTelemetry → Jaeger
│   └── health/                  # Health check handlers standardisés
│
├── python/                      # Packages Python partagés
│   ├── config/                  # Pydantic BaseSettings (toutes les env vars)
│   ├── errors/                  # Exceptions métier (DaedalusError, codes HTTP)
│   ├── middleware/              # FastAPI middleware (auth, CORS, logging)
│   ├── messaging/               # Client aio-pika / RabbitMQ async
│   ├── cache/                   # Client redis.asyncio
│   ├── storage/                 # Client aioboto3 → MinIO
│   ├── tracing/                 # OpenTelemetry setup (FastAPI auto-instrumentation)
│   ├── agents/                  # Utilitaires partagés ADK (tools de base, session factory)
│   └── health/                  # Endpoint /health standardisé
│
├── contracts/                   # Contrats inter-services
│   ├── events/                  # Schémas des events RabbitMQ (JSON Schema / Pydantic)
│   │   ├── user_registered.json
│   │   ├── project_created.json
│   │   ├── procurement_completed.json
│   │   ├── layout_generated.json
│   │   ├── payment_succeeded.json
│   │   └── notification_requested.json
│   ├── openapi/                 # Specs OpenAPI partagées (types communs)
│   └── proto/                   # Définitions gRPC (si utilisé à terme)
│
├── docker/                      # Images Docker de base partagées
│   ├── Dockerfile.python-base   # Image de base Python (ADK, FastAPI, deps)
│   └── Dockerfile.go-base       # Image de base Go (multi-stage)
│
└── k8s/                         # Manifests K8s partagés
    ├── configmaps/              # ConfigMaps communs (URLs Redis, RabbitMQ, etc.)
    ├── secrets-template/        # Templates Secrets (NE PAS committer les valeurs)
    └── rbac/                    # RBAC K8s pour service accounts
```

---

## Variables d'environnement partagées

Tous les services doivent supporter les variables suivantes (injectées via K8s ConfigMap / Secrets) :

### Infrastructure

| Variable | Description | Exemple |
|---|---|---|
| `REDIS_URL` | URI Redis (sentinel ou standalone) | `redis://redis:6379/0` |
| `RABBITMQ_URL` | AMQP URI RabbitMQ | `amqp://user:pass@rabbitmq:5672/` |
| `MINIO_ENDPOINT` | Endpoint MinIO (sans `http://`) | `minio:9000` |
| `MINIO_ACCESS_KEY` | Access key MinIO | `minioadmin` |
| `MINIO_SECRET_KEY` | Secret key MinIO | `minioadmin` |
| `MINIO_BUCKET_3D` | Bucket pour les assets 3D | `daedalus-3d-assets` |
| `MINIO_BUCKET_EXPORTS` | Bucket pour les exports PDF/CSV | `daedalus-exports` |
| `MINIO_USE_SSL` | TLS sur MinIO | `false` (true en prod) |
| `DATABASE_URL` | PostgreSQL DSN | `postgresql+asyncpg://...` |
| `JAEGER_HOST` | Hôte Jaeger (OTLP gRPC) | `jaeger:4317` |
| `SERVICE_NAME` | Nom du service (pour les traces) | `auth-service` |
| `ENVIRONMENT` | Environnement de déploiement | `production` / `development` |
| `LOG_LEVEL` | Niveau de log | `INFO` |

### Auth & Sécurité

| Variable | Description |
|---|---|
| `JWT_SECRET` | Clé de signature JWT (HS256) |
| `JWT_ALGORITHM` | Algorithme JWT | `HS256` |
| `JWT_ACCESS_EXPIRE_MINUTES` | TTL access token (défaut: 30) |
| `JWT_REFRESH_EXPIRE_DAYS` | TTL refresh token (défaut: 7) |
| `OAUTH2_GOOGLE_CLIENT_ID` | Client ID Google OAuth2 |
| `OAUTH2_GOOGLE_CLIENT_SECRET` | Client Secret Google OAuth2 |
| `OAUTH2_GITHUB_CLIENT_ID` | Client ID GitHub OAuth2 |
| `OAUTH2_GITHUB_CLIENT_SECRET` | Client Secret GitHub OAuth2 |
| `TRAEFIK_FORWARD_AUTH_SECRET` | Secret partagé Traefik ↔ auth-service |

### Payment Service

| Variable | Description |
|---|---|
| `PAYSTACK_SECRET_KEY` | Clé secrète Paystack |
| `PAYSTACK_PUBLIC_KEY` | Clé publique Paystack |
| `STRIPE_SECRET_KEY` | Clé secrète Stripe |
| `STRIPE_WEBHOOK_SECRET` | Secret webhook Stripe |
| `PAYMENT_CURRENCY_DEFAULT` | Devise par défaut (`XAF` ou `USD`) |

### Notification Service

| Variable | Description |
|---|---|
| `SMTP_HOST` | Serveur SMTP | `smtp.mailgun.org` |
| `SMTP_PORT` | Port SMTP | `587` |
| `SMTP_USER` | Utilisateur SMTP |
| `SMTP_PASSWORD` | Mot de passe SMTP |
| `SMTP_FROM` | Adresse d'expédition | `no-reply@daedalus.io` |
| `NOTIFICATION_QUEUE` | Queue RabbitMQ pour les notifs | `daedalus.notifications` |

### AI / LLM

| Variable | Description |
|---|---|
| `OPENAI_API_KEY` | Clé OpenAI (GPT-4o) |
| `ANTHROPIC_API_KEY` | Clé Anthropic (Claude) |
| `GOOGLE_GENAI_API_KEY` | Clé Google Gemini (ADK) |
| `MESHY_API_KEY` | Clé Meshy.ai (3D mesh generation) |
| `AGENT_MODEL` | Modèle LLM utilisé par défaut | `gpt-4o` |

---

## Packages Go — Utilisation

### `shared/go/config`
```go
import "github.com/daedalus/shared/go/config"

cfg := config.Load() // Charge depuis l'env
cfg.RedisURL         // string
cfg.RabbitMQURL      // string
```

### `shared/go/errors`
```go
import "github.com/daedalus/shared/go/errors"

return errors.NewNotFound("project", projectID)
return errors.NewUnauthorized("token expired")
return errors.NewValidation("floor dimensions must be positive")
```

### `shared/go/tracing`
```go
import "github.com/daedalus/shared/go/tracing"

shutdown := tracing.Init(ctx, cfg.ServiceName, cfg.JaegerHost)
defer shutdown()
```

### `shared/go/messaging`
```go
import "github.com/daedalus/shared/go/messaging"

publisher := messaging.NewPublisher(cfg.RabbitMQURL)
publisher.Publish(ctx, "daedalus.events", "project.created", payload)
```

---

## Packages Python — Utilisation

### `shared/python/config`
```python
from shared.python.config import Settings

settings = Settings()  # Pydantic lit les env vars automatiquement
settings.redis_url
settings.rabbitmq_url
settings.database_url
```

### `shared/python/errors`
```python
from shared.python.errors import (
    DaedalusNotFoundError,
    DaedalusUnauthorizedError,
    DaedalusValidationError,
    DaedalusPaymentError,
)

raise DaedalusNotFoundError(resource="project", resource_id=project_id)
```

### `shared/python/tracing`
```python
from shared.python.tracing import setup_tracing

# Dans le main.py de chaque service FastAPI
setup_tracing(service_name="auth-service")
```

### `shared/python/messaging`
```python
from shared.python.messaging import RabbitMQClient

client = RabbitMQClient(settings.rabbitmq_url)
await client.publish("daedalus.notifications", {
    "type": "EMAIL",
    "to": "user@example.com",
    "template": "welcome",
    "context": {"username": "Fred"}
})
```

### `shared/python/agents`
```python
from shared.python.agents import create_adk_session, DaedalusTool

# Factory pour créer une session ADK avec config partagée
session = create_adk_session(
    model=settings.agent_model,
    tools=[web_search_tool, minio_upload_tool, ...]
)
```

---

## Contrats d'événements RabbitMQ

Tous les événements inter-services sont publiés sur l'exchange `daedalus.events` de type `topic`.

| Routing Key | Émetteur | Consommateurs |
|---|---|---|
| `user.registered` | auth-service | notification-service |
| `user.verified` | auth-service | notification-service |
| `user.password_reset` | auth-service | notification-service |
| `project.created` | project-service | analytics-service |
| `procurement.started` | orchestrator-agent | procurement-agent |
| `procurement.completed` | procurement-agent | layout-engine, notification-service |
| `layout.generated` | layout-engine | 3d-asset-service, notification-service |
| `asset.ready` | 3d-asset-service | notification-service |
| `payment.succeeded` | payment-service | auth-service, notification-service |
| `payment.failed` | payment-service | notification-service |
| `subscription.created` | payment-service | auth-service |
| `subscription.cancelled` | payment-service | auth-service, notification-service |
| `notification.broadcast` | any service | notification-service |

---

## Conventions de nommage

- **Exchanges RabbitMQ** : `daedalus.events` (topic), `daedalus.notifications` (direct)
- **Queues** : `<service>.<event>` (ex: `notification-service.user.registered`)
- **Buckets MinIO** : `daedalus-3d-assets`, `daedalus-exports`, `daedalus-avatars`
- **Tables Redis** : `cache:<service>:<key>` (ex: `cache:procurement:search:<hash>`)
- **Traces Jaeger** : service name = nom exact du service K8s (ex: `auth-service`)
- **Métriques Prometheus** : préfixe `daedalus_<service>_` (ex: `daedalus_auth_login_total`)

---

## Règles pour les contributeurs

1. **Ne jamais committer** de valeurs de secrets. Utiliser des templates avec `<REPLACE_ME>`.
2. Tout nouveau contrat d'événement **doit être documenté** dans `contracts/events/`.
3. Les packages Go et Python doivent être testés unitairement (couverture ≥ 80 %).
4. Toute config ajoutée doit être documentée dans ce README.
5. Les images Docker de base dans `docker/` doivent rester légères (Alpine/slim).
