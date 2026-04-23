# notification-service

Hub centralisé de toutes les communications sortantes de Daedalus.

## Responsabilités

- **Email transactionnel** : vérification compte, reset password, confirmation paiement, layout généré
- **WebSocket Unicast** : notification temps réel vers un utilisateur spécifique
- **WebSocket Broadcast** : message vers tous les utilisateurs connectés (admin uniquement)
- **Newsletter** : gestion de la liste subscribers (opt-in/out), envoi de campaigns
- **In-app notifications** : persistance des notifs en base (cloche header, lu/non lu)
- **Direct API** : endpoint REST pour envoi direct depuis les services internes

## Stack

- **Python 3.12** + **FastAPI**
- **Socket.io** (python-socketio) pour WebSocket
- **aio-pika** pour RabbitMQ consumer
- **Jinja2** pour les templates email
- **Mailgun API** pour l'envoi SMTP / newsletters

## Architecture interne

```
RabbitMQ Exchange (daedalus.events)
          |
    NotificationConsumer
          |
    DispatchRouter
     /      |      \
Email   WebSocket  In-App DB
(Mailgun) (Socket.io) (PostgreSQL)

Newsletter → Mailgun Lists API
```

## Types de messages supportés

| Type           | Canaux              | Description                              |
|----------------|---------------------|------------------------------------------|
| `INFO`         | WebSocket, In-app   | Information générale                     |
| `SUCCESS`      | Email, WebSocket    | Opération réussie (layout, payment, etc.)|
| `WARNING`      | Email, WebSocket    | Attention requise                        |
| `ERROR`        | Email, WebSocket    | Erreur critique                          |
| `AGENT_UPDATE` | WebSocket           | Mise à jour en temps réel d'un agent IA  |

## Endpoints

```
POST /notifications/send          → Envoi direct (service-to-service, token interne)
GET  /notifications               → Liste des notifs de l'utilisateur connecté
PUT  /notifications/{id}/read     → Marquer comme lue
POST /notifications/newsletter/subscribe   → S'abonner à la newsletter
DELETE /notifications/newsletter/unsubscribe → Se désabonner
WebSocket /ws                     → Connexion WebSocket (auth par JWT)
```

## Events RabbitMQ consommés

| Routing Key              | Action déclenchée                              |
|--------------------------|------------------------------------------------|
| `user.registered`        | Email vérification compte                      |
| `user.verified`          | Email bienvenue                                |
| `user.password_reset`    | Email lien reset password                      |
| `payment.succeeded`      | Email confirmation paiement + upgrade plan     |
| `payment.failed`         | Email échec paiement                           |
| `procurement.completed`  | WebSocket unicast + email résultats sourcing   |
| `layout.generated`       | WebSocket unicast + email layout prêt          |
| `asset.ready`            | WebSocket unicast asset 3D disponible          |
| `notification.broadcast` | WebSocket broadcast tous utilisateurs          |

## Variables d'environnement

```env
RABBITMQ_URL=amqp://...
DATABASE_URL=postgresql+asyncpg://...
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@mg.daedalus.io
SMTP_PASSWORD=...
SMTP_FROM=no-reply@daedalus.io
MAILGUN_API_KEY=...
MAILGUN_DOMAIN=mg.daedalus.io
```
