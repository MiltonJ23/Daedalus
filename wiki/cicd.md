# CI/CD

Chaque service possède son propre `Jenkinsfile` dans `ci-cd/<service>/`.

## Structure

* Étapes communes : checkout · lint · tests unitaires · build image · push registry · deploy.
* Pipelines Go : `go vet`, `go test -race -cover`, build via `Dockerfile`.
* Pipelines Python : `ruff`, `pytest`, build image FastAPI/ADK.
* Pipeline `procurement-agent` : étape supplémentaire Playwright.
* Pipeline `payment-service` : vérification de signature webhook (Paystack/Stripe).
* Pipeline `frontend` : `npm ci`, `npm run lint`, `npm run build`, image Next.js.

## Credentials Jenkins requis

| ID                    | Usage                          |
| --------------------- | ------------------------------ |
| `docker-registry`     | Push images                    |
| `kubeconfig-prod`     | Déploiement Kubernetes         |
| `paystack-webhook`    | Vérif. signature webhook       |
| `sonar-token`         | SonarQube (qualité)            |

Cf. [`../ci-cd/README.md`](../ci-cd/README.md).
