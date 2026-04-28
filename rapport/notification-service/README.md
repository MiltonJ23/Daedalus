# Rapport — notification-service

> SRS v2.0 §3.3 (NOTIF-001 → 009) — **NOUVEAU v2.0**.

## Mission
Hub centralisé sortant : email transactionnel (SMTP/Mailgun), WebSocket
unicast/broadcast, newsletter (Mailgun Lists), notifications in-app
persistées, préférences utilisateur, consumer RabbitMQ.

## Canaux
| Canal       | Use case                                 | Tech                      |
|-------------|------------------------------------------|---------------------------|
| email       | vérification, reset, paiement, layout    | SMTP / Mailgun + Jinja2   |
| websocket   | unicast (« layout prêt »), broadcast     | Socket.io                 |
| in-app      | cloche header, historique non-lu/lu      | PostgreSQL                |
| newsletter  | opt-in/out, tracking ouvertures/clics    | Mailgun Lists API         |

## Types de message (NOTIF-006)
`INFO`, `SUCCESS`, `WARNING`, `ERROR`, `AGENT_UPDATE` — chacun a un template
email distinct.

## Architecture interne
```
RabbitMQ exchange daedalus.events
        ↓
NotificationConsumer
        ↓
DispatchRouter ──┬─ Email (SMTP/Mailgun + Jinja2)
                 ├─ WebSocket (Socket.io)
                 ├─ In-App (PostgreSQL)
                 └─ Newsletter (Mailgun Lists)
```

## API directe (NOTIF-008)
`POST /notifications/send` (token service-to-service `X-Service-Token`) — utilisé
par tout autre service qui ne veut pas passer par RabbitMQ.

## Tests
Mocks SMTP (`aiosmtpd`) + RabbitMQ (`pytest-rabbitmq`) + WebSocket (`websockets.client`).

## Métriques exposées (`/metrics`)
- `notifications_sent_total{channel,type,result}`
- `email_delivery_duration_seconds`
- `websocket_active_connections` (gauge)
- `rabbitmq_consumed_total{event_type}`
