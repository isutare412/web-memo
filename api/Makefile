ENTITY ?= $(error Please set ENTITY)

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: run
run: ## Run server locally.
	go run ./cmd/...

.PHONY: test
test: ## Run tests.
	go test ./...

##@ Entity

.PHONY: ent-new
ent-new: ## Create new entity.
	go run -mod=mod entgo.io/ent/cmd/ent new --target ./internal/core/ent/schema $(ENTITY)

.PHONY: ent-gen
ent-gen: ## Generate Go code for entities.
	go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./internal/core/ent/schema

.PHONY: ent-desc
ent-desc: ## Describe schemas.
	go run -mod=mod entgo.io/ent/cmd/ent describe ./internal/core/ent/schema
