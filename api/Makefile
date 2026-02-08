ENTITY ?= $(error Please set ENTITY)

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: run
run: ## Run server locally.
	go run ./cmd/...

.PHONY: check
check: format lint test ## Run format, lint, and tests

.PHONY: test
test: ginkgo ## Run tests.
	$(GINKGO) --coverprofile coverage.out ./... && go tool cover -func coverage.out > coverage-report.txt

.PHONY: format
format: ## Check Go code formatting
	@echo "🔍 Checking Go code formatting..." && \
	OUTPUT=$$(gofmt -d $$(find . -name '*.go' -not -path './vendor/*')); \
	if [ "$$OUTPUT" ]; then \
		echo "$$OUTPUT"; \
		echo "❌ Code is not properly formatted. Run 'go fmt ./...' to fix."; \
		exit 1; \
	else \
		echo "✅ All Go files are properly formatted."; \
	fi

.PHONY: lint
lint: golangci-lint ## Run golangci-lint
	@echo "🔍 Running golangci-lint..." && \
		$(GOLANGCI_LINT) run && \
		echo "✅ Linting complete."

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
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

.PHONY: mockery
mockery: $(MOCKERY) ## Install mockery if necessary.
$(MOCKERY): $(LOCALBIN)
	@test -s $(MOCKERY) || \
	GOBIN=$(LOCALBIN) go install github.com/vektra/mockery/v2@latest

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Install ginkgo.
$(GINKGO): $(LOCALBIN)
	@test -s $(GINKGO) || GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@v2.20.2


.PHONY: golangci-lint
golangci-lint: ## Check and install golangci-lint if needed
	@GOPATH_BIN=$$(go env GOPATH)/bin; \
	if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		echo "⚠️ golangci-lint not found in $$LOCALBIN. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(LOCALBIN) v2.8.0 && \
		echo "✅ golangci-lint installed successfully."; \
	else \
		echo "✅ golangci-lint is already installed."; \
	fi
