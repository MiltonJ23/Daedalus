# payment-service

Service de gestion des abonnements et paiements pour la plateforme Daedalus.

## Responsabilités

- Gestion des abonnements SaaS (Free / Starter / Business / Enterprise)
- Intégration **Paystack** (cartes + Mobile Money MTN/Orange — marché africain)
- Intégration **Stripe** (cartes internationales VISA/Mastercard)
- Traitement des webhooks paiement avec validation HMAC
- Mise à jour du plan utilisateur après paiement réussi
- Publication d'événements `payment.*` sur RabbitMQ
- Historique des transactions

## Stack

- **Python 3.12** + **FastAPI**
- **Paystack SDK** (paystack-python) + **Stripe SDK**
- **SQLAlchemy async** + PostgreSQL
- **aio-pika** pour RabbitMQ
- **redis.asyncio** pour vérification quotas

## Plans & Pricing

| Plan       | Prix mensuel | Prix annuel (-20%) | Projets | Runs IA/mois | Membres |
|------------|-------------|-------------------|---------|-------------|---------|
| Free       | 0 XAF       | 0 XAF             | 1       | 3           | 1       |
| Starter    | 9 900 XAF   | 95 040 XAF        | 5       | 30          | 1       |
| Business   | 29 900 XAF  | 287 040 XAF       | 20      | 150         | 5       |
| Enterprise | Sur devis   | Négocié           | ∞       | ∞           | ∞       |

## Endpoints

```
POST /payments/initialize      → Créer une transaction Paystack/Stripe
POST /payments/webhook/paystack → Webhook Paystack (HMAC validé)
POST /payments/webhook/stripe   → Webhook Stripe (signature validée)
GET  /payments/plans            → Liste des plans disponibles avec prix
GET  /payments/subscription     → Abonnement courant de l'utilisateur
POST /payments/subscription/cancel → Annuler l'abonnement
GET  /payments/history          → Historique des transactions
POST /payments/upgrade          → Upgrader vers un plan supérieur
```

## Variables d'environnement requises

```env
PAYSTACK_SECRET_KEY=sk_live_...
PAYSTACK_PUBLIC_KEY=pk_live_...
STRIPE_SECRET_KEY=sk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
PAYMENT_CURRENCY_DEFAULT=XAF
DATABASE_URL=postgresql+asyncpg://...
RABBITMQ_URL=amqp://...
REDIS_URL=redis://...
```

## Flux Paystack (Mobile Money)

```
User → POST /payments/initialize
     → payment-service initialise transaction Paystack
     → Retourne authorization_url
     → Frontend redirige vers Paystack Checkout
     → User paie via MTN MoMo / Orange Money / Carte
     → Paystack POST webhook → /payments/webhook/paystack
     → payment-service valide HMAC signature
     → Met à jour user.plan en base
     → Publie payment.succeeded sur RabbitMQ
     → auth-service met à jour X-User-Plan
     → notification-service envoie email confirmation
```

## Events RabbitMQ publiés

| Routing Key           | Payload                                          |
|-----------------------|--------------------------------------------------|
| `payment.succeeded`   | `{user_id, plan, amount, currency, provider}`    |
| `payment.failed`      | `{user_id, plan, reason, provider}`              |
| `subscription.created`| `{user_id, plan, start_date, end_date}`          |
| `subscription.cancelled`| `{user_id, previous_plan, effective_date}`     |

## Impact sur les autres services

- **auth-service** : lit `user.plan` pour inclure `X-User-Plan` dans le Forward Auth header
- **project-service** : vérifie `X-User-Plan` pour limiter les projets actifs
- **orchestrator-agent** : vérifie les quotas IA avant de lancer une session
- **3d-asset-service** : sélectionne la qualité du mesh selon le plan (Basic/Standard/Premium)
- **notification-service** : consomme `payment.*` pour les emails transactionnels

## Mode Partenaire Intégrateur

Les comptes Enterprise peuvent activer le mode Partenaire :
- Création de sous-workspaces clients
- Facturation centralisée (le partenaire paie pour tous ses clients)
- Commission 15% si ≥ 5 clients actifs amenés
- White-label : rapports PDF sans branding Daedalus
- API dédiée pour intégration dans outils métier
