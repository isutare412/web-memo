##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: check
check: format lint test ## Run format, lint, and tests

.PHONY: test
test: ## Run tests.
	go test -race -coverprofile coverage.out ./...

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

.PHONY: generate
generate: ## Generate the OpenAPI client from the api module's spec.
	go generate ./internal/apiclient/...

##@ Build

.PHONY: build
build: ## Build the webmemo binary into bin/.
	go build -o bin/webmemo ./cmd/webmemo

.PHONY: install
install: ## Install the webmemo binary into GOBIN (default ~/go/bin).
	go install ./cmd/webmemo

##@ Tool Dependencies

## Location to install dependencies to.
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

.PHONY: golangci-lint
golangci-lint: ## Check and install golangci-lint if needed
	if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		echo "⚠️ golangci-lint not found in $$LOCALBIN. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(LOCALBIN) v2.11.3 && \
		echo "✅ golangci-lint installed successfully."; \
	else \
		echo "✅ golangci-lint is already installed."; \
	fi
