GO := CGO_ENABLED=0 GO111MODULE=on go

all: build

generate: apiv1 signalcd/proto/agent.pb.go

.PHONY: apiv1
apiv1: api/v1/client api/v1/models api/v1/restapi ui/lib/src/api

SWAGGER ?= docker run --rm \
		--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
		-v $(shell pwd):$(shell pwd) \
		openapitools/openapi-generator-cli:v4.2.3

ui/lib/src/api: signalcd/proto/ui.swagger.json
	-rm -rf ui/lib/src/api
	$(SWAGGER) generate -i $(shell pwd)/signalcd/proto/ui.swagger.json -g dart -o $(shell pwd)/tmp/dart
	mv tmp/dart/lib ui/lib/src/api
	-rm -rf tmp/

signalcd/proto: signalcd/proto/agent.pb.go signalcd/proto/types.pb.go signalcd/proto/ui.pb.go

signalcd/proto/agent.pb.go: signalcd/proto/agent.proto
	protoc signalcd/proto/agent.proto --go_out=plugins=grpc:. -I=. -I=signalcd/proto/vendor

signalcd/proto/types.pb.go: signalcd/proto/types.proto
	protoc signalcd/proto/types.proto --go_out=plugins=grpc:. -I=. -I=signalcd/proto/vendor

signalcd/proto/ui.pb.go: signalcd/proto/ui.proto
	protoc signalcd/proto/ui.proto \
		--go_out=plugins=grpc:. \
		--swagger_out=logtostderr=true:. \
		--grpc-gateway_out=logtostderr=true:. \
		-I=. -I=signalcd/proto/vendor -I=$(GOPATH)/src

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

.PHONY: cmd/checks/http/http
cmd/checks/http/http:
	$(GO) build -v -o ./cmd/checks/http/http ./cmd/checks/http

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

test: test-unit

.PHONY: test-unit
test-unit:
	go test -v -race ./...

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
