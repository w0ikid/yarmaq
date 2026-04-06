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
BOOTSTRAP_COMPOSE := $(DOCKER_DIR)/docker-compose.bootstrap.yml
APPS_COMPOSE     := $(DOCKER_DIR)/docker-compose.yml

# Проекты (для разделения в docker ps)
INFRA_PROJECT    := yarmaq-infra
ZITADEL_PROJECT  := yarmaq-zitadel
APP_PROJECT      := yarmaq-app
APP_SERVICES     := accounts-service transaction-service notification-service

# Цвета для вывода (чтобы логи были читаемыми)
YELLOW := \033[0;33m
NC     := \033[0m

.PHONY: help infra-up infra-down zitadel-up zitadel-down zitadel-bootstrap \
        apps-up apps-down apps-logs app-service-up app-service-start app-service-stop app-service-restart app-service-logs \
        accounts-up accounts-start accounts-stop accounts-restart accounts-logs \
        transactions-up transactions-start transactions-stop transactions-restart transactions-logs \
        notifications-up notifications-start notifications-stop notifications-restart notifications-logs \
        migrate-accounts migrate-transactions migrate-notifications migrate-all \
        up down down-v

# Вспомогательная функция для docker compose с env-файлом
DOCKER_COMPOSE := docker compose --env-file $(ENV_FILE)

define require_app_service
	@if [ -z "$(SERVICE)" ]; then \
		echo "$(YELLOW)SERVICE is required. Example: make $(1) SERVICE=$(word 1,$(APP_SERVICES))$(NC)"; \
		exit 1; \
	fi
	@if ! printf '%s\n' $(APP_SERVICES) | grep -qx -- "$(SERVICE)"; then \
		echo "$(YELLOW)Unknown app service '$(SERVICE)'. Allowed: $(APP_SERVICES)$(NC)"; \
		exit 1; \
	fi
endef

help:
	@echo "$(YELLOW)Yarmaq Makefile commands:$(NC)"
	@echo "  make infra-up          - Запустить Postgres, Kafka и т.д."
	@echo "  make zitadel-up        - Запустить Zitadel"
	@echo "  make zitadel-bootstrap - Выполнить bootstrap Zitadel после zitadel-up"
	@echo "  make apps-up           - Собрать и запустить микросервисы (Docker)"
	@echo "  make app-service-up SERVICE=accounts-service      - Собрать и поднять один микросервис"
	@echo "  make app-service-start SERVICE=accounts-service   - Запустить остановленный микросервис"
	@echo "  make app-service-stop SERVICE=accounts-service    - Остановить один микросервис"
	@echo "  make app-service-restart SERVICE=accounts-service - Перезапустить один микросервис"
	@echo "  make app-service-logs SERVICE=accounts-service    - Логи одного микросервиса"
	@echo "  make accounts-up|start|stop|restart|logs          - Команды для accounts-service"
	@echo "  make transactions-up|start|stop|restart|logs      - Команды для transaction-service"
	@echo "  make notifications-up|start|stop|restart|logs     - Команды для notification-service"
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

zitadel-down-volumes:
	$(DOCKER_COMPOSE) -p $(ZITADEL_PROJECT) -f $(ZITADEL_COMPOSE) down -v

zitadel-bootstrap:
	$(DOCKER_COMPOSE) -f $(BOOTSTRAP_COMPOSE) up --abort-on-container-exit --remove-orphans
	$(DOCKER_COMPOSE) -f $(BOOTSTRAP_COMPOSE) down

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

app-service-up:
	$(call require_app_service,app-service-up)
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) up -d --build $(SERVICE)

app-service-start:
	$(call require_app_service,app-service-start)
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) start $(SERVICE)

app-service-stop:
	$(call require_app_service,app-service-stop)
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) stop $(SERVICE)

app-service-restart:
	$(call require_app_service,app-service-restart)
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) restart $(SERVICE)

app-service-logs:
	$(call require_app_service,app-service-logs)
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) logs -f $(SERVICE)

accounts-up:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) up -d --build accounts-service

accounts-start:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) start accounts-service

accounts-stop:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) stop accounts-service

accounts-restart:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) restart accounts-service

accounts-logs:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) logs -f accounts-service

transactions-up:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) up -d --build transaction-service

transactions-start:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) start transaction-service

transactions-stop:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) stop transaction-service

transactions-restart:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) restart transaction-service

transactions-logs:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) logs -f transaction-service

notifications-up:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) up -d --build notification-service

notifications-start:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) start notification-service

notifications-stop:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) stop notification-service

notifications-restart:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) restart notification-service

notifications-logs:
	$(DOCKER_COMPOSE) -p $(APP_PROJECT) -f $(APPS_COMPOSE) logs -f notification-service

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
