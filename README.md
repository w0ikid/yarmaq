# Yarmaq

Yarmaq is a payment platform built on a microservices architecture in a Go monorepo. Services communicate via Kafka with guaranteed delivery through the Outbox pattern. Transactions are processed through Saga orchestration.

## Requirements

| Tool           | Version |
|----------------|---------|
| Go             | 1.25.5+ |
| Docker         | 24+     |
| Docker Compose | 2.0+    |
| Goose          | latest  |

Install Goose:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Monorepo Structure

```
yarmaq/
├── apps/
│   ├── accounts-service/     # Account and balance management
│   ├── transaction-service/  # Transactions, Saga, Outbox
│   └── notification-service/ # Email/SMS notifications via Kafka
├── pkg/                      # Shared packages (models, kafka, httpclient, etc.)
├── deployment/               # Docker Compose, init scripts
└── secrets/                  # Zitadel JWT keys
```

## Services

### accounts-service
Manages user accounts. Responsible for account creation, balance updates (debit/credit), and operation history via ledger.

### transaction-service
Creates transactions and orchestrates their execution via Saga. Uses the Outbox pattern for reliable event publishing to Kafka.

### notification-service
Consumes events from Kafka (e.g., `transaction.created`) and sends notifications to users.

**Transaction flow:**
```
POST /transactions
  → save PENDING + outbox event
  ← return transaction to client

Outbox Worker
  → read event → publish to Kafka

Consumer "transaction.created"
  → Saga: HOLD (debit) → DEPOSITING (credit) → COMPLETED
  → Notification: Send status update to user
  → on error: REFUND → FAILED
```

## Quick Start

### 1. Configure environment

Create a `.env` file in the project root (or copy `.env.example`):

```env
# ZITADEL
ZITADEL_DOMAIN=zitadel.localhost
ZITADEL_MASTERKEY=MasterkeyNeedsToHave32Characters
ZITADEL_KEY_PATH=secrets/zitadel.json

# POSTGRES (Global for infra)
POSTGRES_USER=danial
POSTGRES_PASSWORD=yarmaq_pass

# SERVICES DB URLs (for Goose)
ACCOUNTS_DB_URL=postgres://danial:yarmaq_pass@localhost:5433/yarmaq_accounts?sslmode=disable
TRANSACTION_DB_URL=postgres://danial:yarmaq_pass@localhost:5433/yarmaq_transactions?sslmode=disable
NOTIFICATION_DB_URL=postgres://danial:yarmaq_pass@localhost:5433/yarmaq_notifications?sslmode=disable
```

### 2. Run the entire stack

The easiest way to start everything (Postgres, Kafka, Zitadel, Migrations, and Microservices) is:

```bash
make up
```

Wait approximately 10-15 seconds for Zitadel and Kafka to initialize.

### 3. Check logs

```bash
make apps-logs
```

## Makefile Commands

| Command                     | Description                                      |
|-----------------------------|--------------------------------------------------|
| `make up`                   | Start everything (infra + apps + migrations)     |
| `make down`                 | Stop everything (keep data)                      |
| `make down-v`               | Stop everything and DELETE volumes               |
| `make infra-up`             | Start Postgres and Kafka only                    |
| `make zitadel-up`           | Start Zitadel only                               |
| `make apps-up`              | Build and start microservices in Docker          |
| `make apps-logs`            | View microservices logs                          |
| `make migrate-all`          | Run migrations for all services                  |
| `make run-accounts`         | Run accounts-service locally (no Docker)         |
| `make run-transactions`     | Run transaction-service locally (no Docker)      |
| `make run-notifications`    | Run notification-service locally (no Docker)     |

## API

### transaction-service (`localhost:8082`)

```
POST /api/v1/transactions
Headers: Authorization: Bearer <token>
         Idempotency-Key: <uuid>

Body:
{
  "from_account_id": "uuid",
  "to_account_id": "uuid",
  "amount": 500,
  "currency": "KZT"
}

GET /api/v1/transactions/:id
```

### accounts-service (`localhost:8081`)

```
POST /api/v1/accounts
GET  /api/v1/accounts/:id
POST /api/v1/accounts/:id/balance
```

## Authentication

[Zitadel](https://zitadel.com) is used as the Identity Provider. JWT tokens are verified via the JWKS endpoint.

- Local domain: `zitadel.localhost:8080`
- Configured via `secrets/zitadel.json`
- Default version: `v2.x` / `v2.42` (check `docker-compose.zitadel.yml`)

## Tech Stack

| Component      | Technology                  |
|----------------|-----------------------------|
| Language       | Go 1.25.5                   |
| HTTP Framework | Fiber v2                    |
| Database       | PostgreSQL 16               |
| ORM            | GORM                        |
| Message Broker | Kafka                       |
| Auth           | Zitadel (OIDC/JWT)          |
| Migrations     | Goose                       |
| Deployment     | Docker Compose + Traefik    |
| Patterns       | Outbox, Saga, Microservices |