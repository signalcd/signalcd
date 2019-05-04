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

.PHONY: cmd/plugins/kubernetes-status/kubernetes-status
cmd/plugins/kubernetes-status/kubernetes-status:
	$(GO) build -v -o ./cmd/plugins/kubernetes-status/kubernetes-status ./cmd/plugins/kubernetes-status

container: cmd/agent/agent container-agent cmd/api/api container-api

.PHONY: container-agent
container-agent:
	docker build -t cd-agent ./cmd/agent

.PHONY: container-api
container-api:
	docker build -t cd-api ./cmd/api

.PHONY: container-kubernetes-status
container-kubernetes-status: cmd/plugins/kubernetes-status/kubernetes-status
	docker build -t quay.io/metalmatze/cd:kubernetes-status ./cmd/plugins/kubernetes-status

.PHONY: container-cheese0
container-cheese0:
	docker build -t quay.io/metalmatze/cd:cheese0 ./examples/cheese0

.PHONY: container-cheese1
container-cheese1:
	docker build -t quay.io/metalmatze/cd:cheese1 ./examples/cheese1

.PHONY: container-cheese2
container-cheese2:
	docker build -t quay.io/metalmatze/cd:cheese2 ./examples/cheese2
