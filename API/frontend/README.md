# API — frontend

Le frontend Next.js 14 ne **publie** pas d'API ; il consomme :

| Service                | Préfixe Traefik           |
|------------------------|---------------------------|
| auth-service           | `/api/auth`               |
| project-service        | `/api/projects`           |
| payment-service        | `/api/payments`           |
| notification-service   | `/api/notifications` + WS |
| orchestrator-agent     | `/api/agent`              |
| procurement-agent      | `/api/procurement`        |
| layout-engine          | `/api/layout`             |
| 3d-asset-service       | `/api/assets`             |
| analytics-service      | `/api/analytics`          |

L'URL de base est injectée via `NEXT_PUBLIC_API_URL`
(cf. `frontend/src/lib/api.ts`).
