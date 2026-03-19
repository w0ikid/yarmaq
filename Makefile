# Подгружаем переменные из .env, если он существует
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# ───────────────────────────────────────────
# Конфигурация
# ───────────────────────────────────────────
GOOSE_DRIVER     ?= postgres
DOCKER_DIR       := ./deployment
DB_COMPOSE       := $(DOCKER_DIR)/docker-compose.db.yml
ZITADEL_COMPOSE  := $(DOCKER_DIR)/docker-compose.zitadel.yml

# Проекты (для разделения в docker ps)
INFRA_PROJECT    := yarmaq-infra
ZITADEL_PROJECT  := yarmaq-zitadel
APP_PROJECT      := yarmaq-app

# Цвета для вывода (чтобы логи были читаемыми)
YELLOW := \033[0;33m
NC     := \033[0m

.PHONY: help infra-up infra-down zitadel-up zitadel-down \
        migrate-accounts migrate-transactions migrate-all \
        up down down-v

help:
	@echo "$(YELLOW)Yarmaq Makefile commands:$(NC)"
	@echo "  make infra-up          - Запустить Postgres, Kafka и т.д."
	@echo "  make zitadel-up        - Запустить Zitadel"
	@echo "  make migrate-all       - Накатить миграции для всех сервисов"
	@echo "  make up                - Поднять всё и накатить миграции"
	@echo "  make down              - Остановить всё (сохранить данные)"
	@echo "  make down-v            - Остановить всё и УДАЛИТЬ данные (volumes)"

# ───────────────────────────────────────────
# ИНФРАСТРУКТУРА
# ───────────────────────────────────────────
infra-up:
	docker compose -p $(INFRA_PROJECT) -f $(DB_COMPOSE) up -d

infra-down:
	docker compose -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down

infra-down-v:
	docker compose -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down -v

zitadel-up:
	docker compose -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) up -d

zitadel-down:
	docker compose -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) down

# ───────────────────────────────────────────
# МИГРАЦИИ (Goose)
# ───────────────────────────────────────────
# Важно: переменные БД должны быть в .env (например ACCOUNTS_DB_URL)
migrate-accounts:
	@echo "$(YELLOW)Running migrations for Accounts Service...$(NC)"
	cd apps/accounts-service && \
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING="$(ACCOUNTS_DB_URL)" \
	goose -dir migrations up

migrate-transactions:
	@echo "$(YELLOW)Running migrations for Transaction Service...$(NC)"
	cd apps/transaction-service && goose -dir migrations $(GOOSE_DRIVER) "$(TRANSACTIONS_DB_URL)" up

migrate-all: migrate-accounts migrate-transactions

# ───────────────────────────────────────────
# ВСЁ ВМЕСТЕ
# ───────────────────────────────────────────
up: infra-up zitadel-up
	@echo "$(YELLOW)Waiting for databases to be ready...$(NC)"
	@sleep 5
	$(MAKE) migrate-all
	@echo "$(YELLOW)Yarmaq is up and running!$(NC)"

down:
	docker compose -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down
	docker compose -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) down

down-v:
	docker compose -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down -v
	docker compose -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) down -v

# ───────────────────────────────────────────
# ЛОКАЛЬНЫЙ ЗАПУСК (Development)
# ───────────────────────────────────────────
run-accounts:
	go run apps/accounts-service/cmd/main.go

run-transactions:
	go run apps/transaction-service/cmd/main.go