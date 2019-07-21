GO := CGO_ENABLED=0 GO111MODULE=on go

all: build

.PHONY: apiv1
apiv1: api/v1/models api/v1/restapi cmd/agent/client cmd/agent/models ui/lib/src/api

GOSWAGGER ?= docker run --rm \
	--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
	-v $(shell pwd):/go/src/github.com/signalcd/signalcd \
	-w /go/src/github.com/signalcd/signalcd quay.io/goswagger/swagger:v0.19.0

api/v1/models api/v1/restapi: swagger.yaml
	-rm -r api/v1/{models,restapi}
	$(GOSWAGGER) generate server -f swagger.yaml --exclude-main -A cd --target api/v1

cmd/agent/client cmd/agent/models: swagger.yaml
	-rm -r cmd/agent/{client,models}
	$(GOSWAGGER) generate client -f swagger.yaml --target cmd/agent

SWAGGER ?= docker run --rm \
		--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
		-v $(shell pwd):/local \
		swaggerapi/swagger-codegen-cli:2.4.0

ui/lib/src/api: swagger.yaml
	-rm -rf ui/lib/src/api
	$(SWAGGER) generate -i /local/swagger.yaml -l dart -o /local/tmp/dart
	mv tmp/dart/lib ui/lib/src/api
	-rm -rf tmp/

.PHONY: build
build: cmd/agent/agent cmd/api/api cmd/ui/ui

.PHONY: cmd/agent/agent
cmd/agent/agent:
	$(GO) build -v -o ./cmd/agent/agent ./cmd/agent

.PHONY: cmd/api/api
cmd/api/api:
	$(GO) build -v -o ./cmd/api/api ./cmd/api

.PHONY: cmd/ui/ui
cmd/ui/ui:
	$(GO) build -v -o ./cmd/ui/ui ./cmd/ui

.PHONY: cmd/plugins/kubernetes-status/kubernetes-status
cmd/plugins/kubernetes-status/kubernetes-status:
	$(GO) build -v -o ./cmd/plugins/kubernetes-status/kubernetes-status ./cmd/plugins/kubernetes-status

.PHONY: ui
ui:
	cd ui && webdev build

.PHONY: ui-serve
ui-serve:
	cd ui && webdev serve

container: container-agent container-api container-kubernetes-status

.PHONY: container-agent
container-agent: cmd/agent/agent
	docker build -t cd-agent ./cmd/agent

.PHONY: container-api
container-api: cmd/api/api
	docker build -t cd-api ./cmd/api

.PHONY: container-ui
container-ui: ui cmd/ui/ui
	cp ui/build/{bulma.min.css,index.html,main.dart.js} ./cmd/ui/assets/
	docker build -t cd-ui ./cmd/ui

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
