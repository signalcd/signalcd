#!/usr/bin/env bash

set -euo pipefail

trap 'kill $(jobs -p); exit 0' EXIT

(
  ./cmd/api/api \
      --ui.assets ./ui
) &

(
  ./cmd/agent/agent \
      --api.url 127.0.0.1:6661 \
      --kubeconfig ~/.kube/config \
      --name 'local' \
      --namespace signalcd-demo \
      --serviceaccount=signalcd-agent
) &

for i in `jobs -p`; do wait $i; done
