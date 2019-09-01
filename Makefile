GO := CGO_ENABLED=0 GO111MODULE=on go

all: build

.PHONY: apiv1
apiv1: api/v1/client api/v1/models api/v1/restapi ui/lib/src/api

GOSWAGGER ?= docker run --rm \
	--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
	-v $(shell pwd):/go/src/github.com/signalcd/signalcd \
	-w /go/src/github.com/signalcd/signalcd quay.io/goswagger/swagger:v0.19.0

api/v1/client api/v1/models api/v1/restapi: swagger.yaml
	-rm -r api/v1/{models,restapi}
	$(GOSWAGGER) generate server -f swagger.yaml --exclude-main -A cd --target api/v1
	$(GOSWAGGER) generate client -f swagger.yaml --target api/v1

SWAGGER ?= docker run --rm \
		--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
		-v $(shell pwd):/local \
		swaggerapi/swagger-codegen-cli:2.4.0

ui/lib/src/api: swagger.yaml
	-rm -rf ui/lib/src/api
	$(SWAGGER) generate -i /local/swagger.yaml -l dart -o /local/tmp/dart
	mv tmp/dart/lib ui/lib/src/api
	-rm -rf tmp/

signalcd/proto/agent.pb.go: signalcd/proto/agent.proto
	protoc signalcd/proto/agent.proto --go_out=plugins=grpc:.

.PHONY: build
build: \
	cmd/agent/agent \
	cmd/api/api \
	cmd/ui/ui \
	cmd/checks/kubernetes-status/kubernetes-status \
	cmd/plugins/drone/drone

.PHONY: cmd/agent/agent
cmd/agent/agent:
	$(GO) build -v -o ./cmd/agent/agent ./cmd/agent

.PHONY: cmd/api/api
cmd/api/api:
	$(GO) build -v -o ./cmd/api/api ./cmd/api

.PHONY: cmd/ui/ui
cmd/ui/ui:
	$(GO) build -v -o ./cmd/ui/ui ./cmd/ui

.PHONY: cmd/checks/kubernetes-status/kubernetes-status
cmd/checks/kubernetes-status/kubernetes-status:
	$(GO) build -v -o ./cmd/checks/kubernetes-status/kubernetes-status ./cmd/checks/kubernetes-status

.PHONY: cmd/plugins/drone/drone
cmd/plugins/drone/drone:
	$(GO) build -v -o ./cmd/plugins/drone/drone ./cmd/plugins/drone

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
container-kubernetes-status: cmd/checks/kubernetes-status/kubernetes-status
	docker build -t quay.io/metalmatze/cd:kubernetes-status ./cmd/checks/kubernetes-status

.PHONY: container-cheese0
container-cheese0:
	docker build -t quay.io/metalmatze/cd:cheese0 ./examples/cheese0

.PHONY: container-cheese1
container-cheese1:
	docker build -t quay.io/metalmatze/cd:cheese1 ./examples/cheese1

.PHONY: container-cheese2
container-cheese2:
	docker build -t quay.io/metalmatze/cd:cheese2 ./examples/cheese2
