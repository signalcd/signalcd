# SignalCD [![Build Status](https://cloud.drone.io/api/badges/signalcd/signalcd/status.svg)](https://cloud.drone.io/signalcd/signalcd)

Continuous Delivery for Kubernetes reacting to Observability Signals.

## Overview

Deploying applications on Kubernetes often involves a lot more manual steps then we want.  
We want to reuse existing observability signals to automate all steps in the application lifecycle.
This will drastically reduce the chance of human errors when deploying business critical applications.

## Example

```yaml
kind: Delivery
name: example

steps:
  name: deploy
  image: quay.io/signalcd/example:$DRONE_COMMIT
  commands:
    - kubectl apply -f /manifests

checks:
  name: kubernetes-status
  image: quay.io/signalcd/kubernetes-status
  labels: app=cheese
  duration: 60
```
