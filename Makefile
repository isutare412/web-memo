ifneq (,$(wildcard ./.env.local))
    include .env.local
    export $(shell sed 's/=.*//' ./.env.local)
endif

# images
TAG_API ?= latest
TAG_UI ?= latest
IMAGE_API = redshoore/webmemo-api:$(TAG_API)
IMAGE_UI = redshoore/webmemo-ui:$(TAG_UI)

# docker hub
DOCKER_USER ?= <docker_hub_username>
DOCKER_PASSWORD ?= <docker_hub_secret>

# docker compose
TARGET ?=
ENV_FILE ?= .env
COMPOSE_CMD = docker compose -f compose.yaml --env-file $(ENV_FILE)

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Protocol Buffers

.PHONY: proto-go
proto-go: protoc protoc-gen-go protoc-gen-go-grpc ## Compile protobuf to Go.
	@PATH=$(LOCALPATH) protoc \
		-Iprotos \
		--go_out=./api/internal \
		--go_opt=paths=import \
		--go-grpc_out=./api/internal \
		--go-grpc_opt=paths=import \
		protos/webmemo/*.proto

##@ Build

.PHONY: build-api
build-api: ## Build docker image of API.
	docker build -f api/Dockerfile -t $(IMAGE_API) api

.PHONY: build-ui
build-ui: ## Build docker image of UI.
	docker build -f ui/Dockerfile -t $(IMAGE_UI) ui

.PHONY: push-api
push-api: ## Push docker image of API.
	echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USER) --password-stdin
	docker push $(IMAGE_API)

.PHONY: push-ui
push-ui: ## Push docker image of UI.
	echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USER) --password-stdin
	docker push $(IMAGE_UI)

##@ Docker Compose

.PHONY: up
up: ## Run components.
	$(COMPOSE_CMD) up -d $(TARGET)

.PHONY: down
down: ## Shutdown components.
	$(COMPOSE_CMD) down $(TARGET)

.PHONY: ps
ps: ## Print running components.
	$(COMPOSE_CMD) ps $(TARGET)

.PHONY: logs
logs: ## Tail logs of components.
	$(COMPOSE_CMD) logs -f $(TARGET)

##@ Tools

# Location to install dependencies to.
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# Modified path environment variable including dependencies.
LOCALPATH ?= $(LOCALBIN):$(PATH)

# Tool Binaries
PROTOC_GEN_GO ?= $(LOCALBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC ?= $(LOCALBIN)/protoc-gen-go-grpc

# Tool Versions
PROTOC_GEN_GO_VERSION ?= v1.34.2
PROTOC_GEN_GO_GRPC_VERSION ?= v1.4.0

.PHONY: protoc
protoc: ## Install protoc.
	@command -v protoc > /dev/null || brew install protoc

.PHONY: protoc-gen-go
protoc-gen-go: $(PROTOC_GEN_GO) ## Install protoc-gen-go locally if necessary.
$(PROTOC_GEN_GO): $(LOCALBIN)
	@test -s $(LOCALBIN)/protoc-gen-go || \
	GOBIN=$(LOCALBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

.PHONY: protoc-gen-go-grpc
protoc-gen-go-grpc: $(PROTOC_GEN_GO_GRPC) ## Install protoc-gen-go-grpc locally if necessary.
$(PROTOC_GEN_GO_GRPC): $(LOCALBIN)
	@test -s $(LOCALBIN)/protoc-gen-go-grpc || \
	GOBIN=$(LOCALBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
