# Architecture

Daedalus suit une architecture microservices conforme à la SRS v2.0.

## Vue d'ensemble

* **Frontend** Next.js 14 + React Three Fiber (visualisation 3D).
* **API Gateway** Traefik avec ForwardAuth vers `auth-service`.
* **Bus d'événements** RabbitMQ — exchange `daedalus.events`.
* **Cache & quotas** Redis — clés `plan:quota:<user_id>:runs`.
* **Object storage** MinIO (compatible S3) — assets 3D, exports.
* **Persistance** PostgreSQL 16 (par service).

## Services

| Service                | Stack             | Port |
| ---------------------- | ----------------- | ---- |
| project-service        | Go + pgx          | 8080 |
| auth-service           | Go + Fiber        | 8081 |
| payment-service        | Python + FastAPI  | 8082 |
| notification-service   | Python + FastAPI  | 8083 |
| orchestrator-agent     | Python + ADK      | 8084 |
| procurement-agent      | Python + ADK      | 8085 |
| layout-engine          | Python + ADK      | 8086 |
| 3d-asset-service       | Python + Meshy    | 8087 |
| analytics-service      | Python + FastAPI  | 8088 |

## Diagrammes

Voir [`../architecture/project-management-service/diagrams/`](../architecture/project-management-service/diagrams/).
