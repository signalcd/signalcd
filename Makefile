GO := CGO_ENABLED=0 GO111MODULE=on go

all: build

run: build
	./development/run.sh

# Code Generation

generate: api

api: api/client/go api/client/javascript api/server/go

OPENAPI ?= docker run --rm \
		--user=$(shell id -u $(USER)):$(shell id -g $(USER)) \
		-v $(shell pwd):$(shell pwd) \
		openapitools/openapi-generator-cli:v4.3.1

GOIMPORTS ?= goimports

api/client/go: api/api.yaml
	-rm -rf $@
	$(OPENAPI) generate -i $(shell pwd)/api/api.yaml -g go -o $(shell pwd)/api/client/go --additional-properties=withGoCodegenComment=true
	-rm -rf $@/go.mod
	-rm -rf $@/go.sum
	$(GOIMPORTS) -w $(shell find ./api/client/go/ -name '*.go')
	touch $@

api/client/javascript: api/api.yaml
	-rm -rf $@
	$(OPENAPI) generate -i $(shell pwd)/api/api.yaml -g javascript -o $(shell pwd)/api/client/javascript --additional-properties=usePromises=true
	touch $@

api/server/go: api/api.yaml
	-rm -rf $@
	$(OPENAPI) generate -i $(shell pwd)/api/api.yaml -g go-server -o $(shell pwd)/api/server/go
	-rm -rf $@/{go.mod,main.go}
	$(GOIMPORTS) -w $(shell find ./api/server/go/ -name '*.go')
	touch $@

# Building Binaries with Go

.PHONY: build
build: \
	cmd/agent/agent \
	cmd/api/api \
	cmd/checks/kubernetes-status/kubernetes-status \
	cmd/plugins/drone/drone \
	ui/bundle.js

.PHONY: cmd/agent/agent
cmd/agent/agent:
	$(GO) build -v -o ./cmd/agent/agent ./cmd/agent

.PHONY: cmd/api/api
cmd/api/api:
	$(GO) build -v -o ./cmd/api/api ./cmd/api

.PHONY: cmd/checks/http/http
cmd/checks/http/http:
	$(GO) build -v -o ./cmd/checks/http/http ./cmd/checks/http

.PHONY: cmd/checks/kubernetes-status/kubernetes-status
cmd/checks/kubernetes-status/kubernetes-status:
	$(GO) build -v -o ./cmd/checks/kubernetes-status/kubernetes-status ./cmd/checks/kubernetes-status

.PHONY: cmd/plugins/drone/drone
cmd/plugins/drone/drone:
	$(GO) build -v -o ./cmd/plugins/drone/drone ./cmd/plugins/drone

test: test-unit

.PHONY: test-unit
test-unit:
	go test -v -race ./...

# Building the UI

node_modules: package.json package-lock.json
	npm install
	touch $@

ui/bundle.js: node_modules $(shell find ./ui -iname '*.js' | grep -v ./ui/bundle.js)
	node_modules/esbuild/bin/esbuild --bundle ui/index.js --outfile=ui/bundle.js --sourcemap

# Building Containers

container: container-agent container-api container-kubernetes-status

.PHONY: container-agent
container-agent: cmd/agent/agent
	docker build -t cd-agent ./cmd/agent

.PHONY: container-api
container-api: cmd/api/api ui/bundle.js
	mkdir -p ./cmd/api/assets
	cp ./ui/index.html ./cmd/api/assets
	cp ./ui/bundle.js ./cmd/api/assets
	docker build -t cd-api ./cmd/api

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
