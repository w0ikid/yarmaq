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
│   └── transaction-service/  # Transactions, Saga, Outbox
├── pkg/                      # Shared packages (models, kafka, httpclient, etc.)
├── deployment/               # Docker Compose, init scripts
└── secrets/                  # Zitadel JWT keys
```

## Services

### accounts-service
Manages user accounts. Responsible for account creation, balance updates (debit/credit), and operation history via ledger.

### transaction-service
Creates transactions and orchestrates their execution via Saga. Uses the Outbox pattern for reliable event publishing to Kafka.

**Transaction flow:**
```
POST /transactions
  → save PENDING + outbox event
  ← return transaction to client

Outbox Worker
  → read event → publish to Kafka

Consumer "transaction.created"
  → Saga: HOLD (debit) → DEPOSITING (credit) → COMPLETED
  → on error: REFUND → FAILED
```

## Quick Start

### 1. Configure environment

Create a `.env` file in the project root (or copy and adjust the example below):

```env
# ACCOUNTS SERVICE
ACCOUNTS_APP_PORT=8081
ACCOUNTS_POSTGRES_DB_NAME=yarmaq_accounts
ACCOUNTS_KAFKA_BROKERS=localhost:9092

# TRANSACTION SERVICE
TRANSACTION_APP_PORT=8082
TRANSACTION_POSTGRES_DB_NAME=yarmaq_transactions

# POSTGRES
POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_USER=danial
POSTGRES_PASSWORD=yarmaq_pass
POSTGRES_SSLMODE=disable

# MIGRATION URLS
ACCOUNTS_DB_URL=postgres://danial:yarmaq_pass@localhost:5433/yarmaq_accounts?sslmode=disable
TRANSACTION_DB_URL=postgres://danial:yarmaq_pass@localhost:5433/yarmaq_transactions?sslmode=disable

# ZITADEL
ZITADEL_DOMAIN=zitadel.localhost
ZITADEL_EXTERNALPORT=8080
ZITADEL_MASTERKEY=MasterkeyNeedsToHave32Characters
ZITADEL_KEY_PATH=secrets/zitadel.json
```

> ⚠️ `ZITADEL_MASTERKEY` must be exactly 32 characters. The values above are for **local development only** — change them before exposing to any network.

### 2. Start infrastructure + run migrations

```bash
make up
```

This starts Postgres, Kafka, Zookeeper, Zitadel and runs all migrations automatically.

### 3. Run services

```bash
# Terminal 1
make run-accounts

# Terminal 2
make run-transactions
```

## Makefile Commands

| Command                     | Description                                 |
|-----------------------------|---------------------------------------------|
| `make up`                   | Start all infrastructure and run migrations |
| `make down`                 | Stop all containers (keep data)             |
| `make down-v`               | Stop all containers and delete volumes      |
| `make infra-up`             | Start Postgres, Kafka only                  |
| `make zitadel-up`           | Start Zitadel only                          |
| `make migrate-all`          | Run migrations for all services             |
| `make migrate-accounts`     | Run migrations for accounts-service         |
| `make migrate-transactions` | Run migrations for transaction-service      |
| `make run-accounts`         | Run accounts-service locally                |
| `make run-transactions`     | Run transaction-service locally             |

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
Headers: Authorization: Bearer <token>
```

### accounts-service (`localhost:8081`)

```
POST /api/v1/accounts
GET  /api/v1/accounts/:id
POST /api/v1/accounts/:id/balance
```

## Authentication

[Zitadel](https://zitadel.com) is used as the Identity Provider. JWT tokens are verified via the JWKS endpoint.

- Zitadel runs on `http://zitadel.localhost:8080`
- Place the service key in `secrets/zitadel.json`
- Default version: `v4.11.0`

## Tech Stack

| Component      | Technology                 |
|----------------|----------------------------|
| Language       | Go 1.25.5                  |
| HTTP           | Fiber v2                   |
| ORM            | GORM + PostgreSQL 16       |
| Message Broker | Kafka (segmentio/kafka-go) |
| Auth           | Zitadel v4.11.0 (OIDC/JWT) |
| Migrations     | Goose                      |
| Patterns       | Outbox, Saga, CQRS         |