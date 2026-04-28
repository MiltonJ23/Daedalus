# Monitoring

## Prometheus

Configuration : [`infrastructure/monitoring/prometheus.yml`](../infrastructure/monitoring/prometheus.yml).

Tous les services exposent `/metrics`. Le `project-service` instrumente via
`github.com/prometheus/client_golang` :

* `http_requests_total{method,path,status}` — compteur de requêtes.
* `http_request_duration_seconds{method,path}` — histogramme de latence.
* `http_requests_in_flight` — jauge concurrente.
* `http_errors_total{method,path,status}` — erreurs 5xx.
* `project_operations_total{operation,result}` — métriques métier.

## Grafana

Dashboards préconfigurés dans `infrastructure/monitoring/grafana/`.

## Tracing

OpenTelemetry → Jaeger (cf. SRS §5.3, NFR-O02).

## Alerting

Alertmanager + intégration Slack/PagerDuty (à configurer par environnement).
