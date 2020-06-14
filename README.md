# SignalCD [![Build Status](https://cloud.drone.io/api/badges/signalcd/signalcd/status.svg)](https://cloud.drone.io/signalcd/signalcd)

Continuous Delivery for Kubernetes reacting to Observability Signals.

## Overview

Deploying applications on Kubernetes often involves a lot more manual steps then we want.  
We want to reuse existing observability signals to automate all steps in the application lifecycle.
This will drastically reduce the chance of human errors when deploying business critical applications.

## Example

[embedmd]:# "examples/config/01-simple.yaml"
```yaml
kind: Pipeline
name: example

steps:
  - name: deploy
    image: quay.io/signalcd/example
    commands:
      - kubectl apply -f /manifests

checks:
  - name: kubernetes-status
    image: quay.io/signalcd/kubernetes-status
    labels: app=cheese
    duration: 10m
```

## Development
### Architecture

![architecture.svg](documentation/architecture.svg)

### Prerequisites 

* Go 1.14+
* Docker (needed to generate OpenAPI spec and build containers)

### API

The API is generated with an OpenAPI spec that you can find in `/api/api.yaml`. From that file we generate a Go server into `/api/go-server`, a Go client into `/api/go` and a JavaScript client into `/api/javascript`. All of these can be generate by running `make api`.