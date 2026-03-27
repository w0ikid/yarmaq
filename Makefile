# Файл переменных окружения (по умолчанию .env)
ENV_FILE         ?= .env
ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
endif

# ───────────────────────────────────────────
# Конфигурация
# ───────────────────────────────────────────
GOOSE_DRIVER     ?= postgres
DOCKER_DIR       := ./deployment
DB_COMPOSE       := $(DOCKER_DIR)/docker-compose.infra.yml
ZITADEL_COMPOSE  := $(DOCKER_DIR)/docker-compose.zitadel.yml
APPS_COMPOSE     := $(DOCKER_DIR)/docker-compose.yml

# Проекты (для разделения в docker ps)
INFRA_PROJECT    := yarmaq-infra
ZITADEL_PROJECT  := yarmaq-zitadel
APP_PROJECT      := yarmaq-app

# Цвета для вывода (чтобы логи были читаемыми)
YELLOW := \033[0;33m
NC     := \033[0m

.PHONY: help infra-up infra-down zitadel-up zitadel-down \
        apps-up apps-down apps-logs \
        migrate-accounts migrate-transactions migrate-notifications migrate-all \
        up down down-v

# Вспомогательная функция для docker compose с env-файлом
DOCKER_COMPOSE := docker compose --env-file $(ENV_FILE)

help:
	@echo "$(YELLOW)Yarmaq Makefile commands:$(NC)"
	@echo "  make infra-up          - Запустить Postgres, Kafka и т.д."
	@echo "  make zitadel-up        - Запустить Zitadel"
	@echo "  make apps-up           - Собрать и запустить микросервисы (Docker)"
	@echo "  make apps-logs         - Просмотр логов микросервисов"
	@echo "  make migrate-all       - Накатить миграции для всех сервисов"
	@echo "  make up                - Поднять всё (infra + zitadel + migrations + apps)"
	@echo "  make down              - Остановить всё (сохранить данные)"
	@echo "  make down-v            - Остановить всё и УДАЛИТЬ данные (volumes)"

# ───────────────────────────────────────────
# ИНФРАСТРУКТУРА
# ───────────────────────────────────────────
infra-up:
	$(DOCKER_COMPOSE) -p $(INFRA_PROJECT) -f $(DB_COMPOSE) up -d

infra-down:
	$(DOCKER_COMPOSE) -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down

infra-down-v:
	$(DOCKER_COMPOSE) -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down -v

zitadel-up:
	$(DOCKER_COMPOSE) -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) up -d --wait

zitadel-down:
	$(DOCKER_COMPOSE) -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) down

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
	cd apps/transaction-service && \
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING="$(TRANSACTION_DB_URL)" \
	goose -dir migrations up

migrate-notifications:
	@echo "$(YELLOW)Running migrations for Notification Service...$(NC)"
	cd apps/notification-service && \
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING="$(NOTIFICATION_DB_URL)" \
	goose -dir migrations up

migrate-all: migrate-accounts migrate-transactions migrate-notifications

# ───────────────────────────────────────────
# ПРИЛОЖЕНИЯ (Microservices)
# ───────────────────────────────────────────
apps-up:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) up -d --build

apps-down:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) down -v

apps-logs:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) logs -f

# ───────────────────────────────────────────
# ВСЁ ВМЕСТЕ
# ───────────────────────────────────────────
up: infra-up zitadel-up
	@echo "$(YELLOW)Waiting for infra to be ready...$(NC)"
	@sleep 10
	$(MAKE) ENV_FILE=$(ENV_FILE) migrate-all
	$(MAKE) ENV_FILE=$(ENV_FILE) apps-up
	@echo "$(YELLOW)Yarmaq is up and running!$(NC)"

down:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) down
	$(DOCKER_COMPOSE) -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down
	$(DOCKER_COMPOSE) -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) down

down-v:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) down -v
	$(DOCKER_COMPOSE) -p $(INFRA_PROJECT) -f $(DB_COMPOSE) down -v
	$(DOCKER_COMPOSE) -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) down -v

# ───────────────────────────────────────────
# ЛОКАЛЬНЫЙ ЗАПУСК (Development)
# ───────────────────────────────────────────
run-accounts:
	go run apps/accounts-service/cmd/api/main.go

run-transactions:
	go run apps/transaction-service/cmd/api/main.go

run-notifications:
	go run apps/notification-service/cmd/api/main.go