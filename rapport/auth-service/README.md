# Rapport — auth-service

> SRS v2.0 §3.1 — Authentification & autorisation.

## Mission
Source unique d'identité Daedalus : enregistrement classique, OAuth2.0
Google/GitHub, JWT, refresh, RBAC, et **endpoint Forward Auth** consommé par
Traefik avant tout routage vers les services protégés.

## Endpoints clés
- `POST /auth/register`, `POST /auth/login` (AUTH-001, AUTH-002)
- `GET  /auth/oauth/{google,github}` + `/callback` (AUTH-005, AUTH-006)
- `GET  /auth/verify` → headers `X-User-Id|Role|Plan` (AUTH-007, AUTH-008)
- `POST /auth/password/reset` (AUTH-003)

## Stack
Go 1.22 + Fiber, PostgreSQL (users, oauth_tokens), Redis (sessions OAuth, TTL
aligné), bcrypt, JWT (HS256, secret rotation envisagée).

## Politique tokens
| Token   | TTL    | Stockage                    |
|---------|--------|-----------------------------|
| Access  | 30 min | Header `Authorization`      |
| Refresh | 7 j    | Cookie HttpOnly + DB        |
| OAuth   | aligné | Redis chiffré (NFR-SEC02)   |

## Tests requis
- ≥ 80% couverture (SRS Sprint 2 §6.2).
- Cas : login OK/KO, refresh, JWT expiré, OAuth callback, RBAC 403, ForwardAuth.

## Métriques exposées (`/metrics`)
- `http_requests_total{path,method,status}`
- `http_request_duration_seconds{path}`
- `auth_login_total{result}`
- `auth_oauth_login_total{provider,result}`
- `forwardauth_verify_duration_seconds`
