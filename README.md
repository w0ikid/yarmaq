# Yarmaq

Yarmaq is a payment platform built on a microservices architecture in a Go monorepo. Services communicate via Kafka with guaranteed delivery through the Outbox pattern. Transactions are processed through Saga orchestration.

## Requirements

| Tool | Version |
|------|---------|
| Go | 1.25.5+ |
| Docker | 24+ |
| Docker Compose | 2.0+ |

## Monorepo Structure

```
yarmaq/
├── apps/
│   ├── accounts-service/
│   ├── transaction-service/
│   └── notification-service/
├── pkg/
├── deployment/
└── secrets/
```

## Transaction Flow

```
POST /transactions
  → save PENDING + outbox event
  ← return transaction to client

Outbox Worker
  → read event → publish to Kafka

Consumer "transaction.created"
  → Saga: HOLD → DEPOSIT → COMPLETED
  → Notification: Send status update to user
  → on error: REFUND → FAILED
```

## Quick Start

See [docs/Quickstart.md](docs/Quickstart.md).

All available commands are in the `Makefile`.
If you prefer [Task](https://taskfile.dev), a `Taskfile.yml` is also available.

## Local Services

| Service | URL |
|---------|-----|
| Zitadel Console | http://zitadel.localhost:8080/ui/console |
| Mailhog | http://localhost:8025 |
| Kafka UI | http://localhost:8085 |

## API

Postman collection is available at `docs/schema/yarmaq.postman_collection.json`.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25.5 |
| HTTP Framework | Fiber v2 |
| Database | PostgreSQL 16 |
| Message Broker | Kafka |
| Auth | Zitadel (OIDC/JWT) |
| Migrations | Goose |
| Deployment | Docker Compose + Traefik |
| Patterns | Outbox, Saga, Microservices |