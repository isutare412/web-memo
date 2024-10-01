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
test: ginkgo ## Run tests.
	$(GINKGO) --coverprofile coverage.out ./... && go tool cover -func coverage.out > coverage-report.txt

.PHONY: mocks
mocks: mockery ## Generate mock implementations.
	$(MOCKERY) --config .mockery.yaml

##@ Entity

.PHONY: ent-new
ent-new: ## Create new entity.
	go run -mod=mod entgo.io/ent/cmd/ent new --target ./internal/core/ent/schema $(ENTITY)
	@go mod tidy

.PHONY: ent-gen
ent-gen: ## Generate Go code for entities.
	go run -mod=mod entgo.io/ent/cmd/ent generate --feature=sql/upsert --feature=sql/execquery ./internal/core/ent/schema
	@go mod tidy

.PHONY: ent-desc
ent-desc: ## Describe schemas.
	go run -mod=mod entgo.io/ent/cmd/ent describe ./internal/core/ent/schema
	@go mod tidy

##@ Tool Dependencies

## Location to install dependencies to.
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## Tool Binaries
MOCKERY ?= $(LOCALBIN)/mockery
GINKGO ?= $(LOCALBIN)/ginkgo

.PHONY: mockery
mockery: $(MOCKERY) ## Install mockery if necessary.
$(MOCKERY): $(LOCALBIN)
	@test -s $(MOCKERY) || \
	GOBIN=$(LOCALBIN) go install github.com/vektra/mockery/v2@latest

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Install ginkgo.
$(GINKGO): $(LOCALBIN)
	@test -s $(GINKGO) || GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@v2.20.2
