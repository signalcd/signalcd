GO := CGO_ENABLED=0 GO111MODULE=on go

all: build

.PHONY: build
build: cmd/agent/agent cmd/api/api

.PHONY: cmd/agent/agent
cmd/agent/agent:
	$(GO) build -v -o ./cmd/agent/agent ./cmd/agent

.PHONY: cmd/api/api
cmd/api/api:
	$(GO) build -v -o ./cmd/api/api ./cmd/api

container: cmd/agent/agent container-agent cmd/api/api container-api

.PHONY: container-agent
container-agent:
	docker build -t cd-agent ./cmd/agent

.PHONY: container-api
container-api:
	docker build -t cd-api ./cmd/api
