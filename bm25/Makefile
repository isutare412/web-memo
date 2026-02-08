##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: run
run: ## Run server locally.
	uv run uvicorn main:app --reload --port 28081

.PHONY: check
check: format lint ## Run format and lint checks.

.PHONY: format
format: ## Check Python code formatting.
	@echo "🔍 Checking Python code formatting..." && \
	uv run ruff format --check . && \
	echo "✅ All Python files are properly formatted."

.PHONY: lint
lint: ## Run ruff linter.
	@echo "🔍 Running ruff linter..." && \
	uv run ruff check . && \
	echo "✅ Linting complete."

.PHONY: fix
fix: ## Auto-fix formatting and lint issues.
	uv run ruff format .
	uv run ruff check --fix .

.PHONY: lock
lock: ## Generate uv.lock.
	uv lock
