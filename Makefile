# docker compose
TARGET ?=
ENV_FILE ?= .env
COMPOSE_CMD = docker compose -f compose.yaml --env-file $(ENV_FILE)

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Docker Compose

.PHONY: up
up: ## Run components.
	$(COMPOSE_CMD) up -d $(TARGET)

.PHONY: down
down: ## Shutdown components.
	$(COMPOSE_CMD) down $(TARGET)

.PHONY: p
ps: ## Print running components.
	$(COMPOSE_CMD) ps $(TARGET)

.PHONY: log
logs: ## Tail logs of components.
	$(COMPOSE_CMD) logs -f $(TARGET)
