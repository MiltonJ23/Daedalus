# Rapport — payment-service

> SRS v2.0 §3.4 (PAY-001 → 011) — **NOUVEAU v2.0**.

## Mission
Brique SaaS Freemium : 4 plans (Free / Starter / Business / Enterprise),
souscription Paystack (cartes + Mobile Money) et Stripe (cartes intl.),
webhooks signés, gestion upgrade/downgrade, période d'essai 14 j (Starter),
mode Partenaire Intégrateur.

## Grille tarifaire
| Plan       | XAF/mois | USD/mois | Projets | Runs IA/mois | Membres |
|------------|---------:|---------:|--------:|-------------:|--------:|
| Free       |        0 |       $0 |       1 |            3 |       1 |
| Starter    |    9 900 |      $16 |       5 |           30 |       1 |
| Business   |   29 900 |      $49 |      20 |          150 |       5 |
| Enterprise |  Sur dev |   Custom |       ∞ |            ∞ |       ∞ |

Prix annuel = -20% (PAY-003). Prix XAF par défaut, conversion USD temps réel
(PAY-009, ANA-001).

## Flux paiement (annexe SRS §9.2)
1. `POST /payments/initialize` (front)
2. payment-service crée la transaction Paystack/Stripe
3. Front redirige vers la hosted page
4. Webhook signé (HMAC) → `POST /payments/webhook/{paystack|stripe}`
5. Mise à jour plan + publication `payment.succeeded` sur RabbitMQ
6. auth-service met à jour `users.plan`
7. notification-service envoie l'email confirmation
8. Front reçoit l'event WebSocket → confirmation UI

## Sécurité (NFR-SEC03)
- Webhooks validés par signature HMAC SHA512 (Paystack) / Stripe-Signature.
- Idempotence via `transaction_id` + `provider_event_id` UNIQUE.

## Tests
- ≥ 80% couverture, focus signature webhook + transitions d'état (pending →
  succeeded → refunded).

## Métriques exposées (`/metrics`)
- `payments_total{provider,plan,result}`
- `webhook_received_total{provider,event_type,signature_valid}`
- `subscription_changes_total{from,to}`
- `mrr_xaf` (gauge, exporté périodiquement)
