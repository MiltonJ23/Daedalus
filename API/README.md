# API — Daedalus

Un dossier par module/service. Chaque dossier contient au minimum :

- `openapi.yaml` (ou `asyncapi.yaml` pour les services event-driven) — contrat
  d'interface.
- `README.md` — résumé des endpoints, codes d'erreur, exemples.

| Module                  | Type      | Spec                                  |
|-------------------------|-----------|---------------------------------------|
| `auth-service/`         | REST      | OAuth2.0, JWT, /auth/verify (Forward Auth) |
| `project-service/`      | REST      | CRUD projets factory                  |
| `payment-service/`      | REST      | Plans, abonnements, webhooks Paystack/Stripe |
| `notification-service/` | REST + WS | Email, Unicast, Broadcast, Newsletter |
| `orchestrator-agent/`   | REST + SSE | POST /agent/run, streaming logs      |
| `procurement-agent/`    | REST + WS | Sourcing IA, ranking, approval flow   |
| `layout-engine/`        | REST      | POST /layout/generate → manifeste JSON |
| `3d-asset-service/`     | REST      | POST /assets/generate → URL GLB MinIO |
| `analytics-service/`    | REST      | BOM, ROI, exports PDF/CSV             |
| `frontend/`             | —         | Consommateur (pas d'API exposée)      |

Toutes les specs sont versionnées avec le service correspondant et publiées
dans Swagger UI à `/api/docs` (cf. SRS §6.2 Sprint 7).
