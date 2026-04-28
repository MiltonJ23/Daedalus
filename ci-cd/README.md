# CI-CD — Daedalus

Un dossier par module/service. Chaque dossier contient un `Jenkinsfile` autonome.

| Module                  | Tech              | Pipeline                                              |
|-------------------------|-------------------|-------------------------------------------------------|
| `auth-service/`         | Go + Fiber        | go vet → go test → docker → push → k3s rollout        |
| `project-service/`      | Go + pgx/v5       | go vet → go test → docker → push → k3s rollout        |
| `payment-service/`      | Python + FastAPI  | ruff → pytest → webhook sig tests → docker → k3s      |
| `notification-service/` | Python + FastAPI  | ruff → pytest → docker → k3s rollout                  |
| `orchestrator-agent/`   | Python + ADK      | ruff → pytest → docker → k3s rollout                  |
| `procurement-agent/`    | Python + Playwright | playwright install → pytest → docker → k3s          |
| `layout-engine/`        | Python + NetworkX | ruff → pytest → docker → k3s rollout                  |
| `3d-asset-service/`     | Python + Meshy.ai | ruff → pytest → docker → k3s rollout                  |
| `analytics-service/`    | Python + FastAPI  | ruff → pytest → docker → k3s rollout                  |
| `frontend/`             | Next.js 14        | npm ci → lint → next build → docker → k3s rollout     |

## Configuration Jenkins requise

Credentials à créer dans Jenkins :

- `docker-registry-url` — URL du registry (ex. `registry.daedalus.io`)
- `docker-registry-credentials` — username/password
- `kubeconfig` — fichier kubeconfig pour le cluster K3s

## Convention

Chaque pipeline :
1. Se déclenche sur `githubPush()` ;
2. Tag les images avec le SHA court du commit + `latest` ;
3. Déploie via `kubectl set image` + `kubectl rollout status` ;
4. Nettoie le workspace (`cleanWs()`).

> NB : le pipeline monolithique historique est conservé à `ci-cd/jenkins/Jenkinsfile`
> pour rétro-compatibilité.
