#!/usr/bin/env bash

set -euo pipefail

trap 'kill $(jobs -p); exit 0' EXIT

(
  echo "running api"
  ./cmd/api/api \
      --ui.assets ./ui
) &

sleep 1

(
  echo "running agent"
  ./cmd/agent/agent \
      --api.url http://127.0.0.1:6660 \
      --kubeconfig ~/.kube/config \
      --name 'local' \
      --namespace signalcd-demo \
      --serviceaccount=signalcd-agent
) &

for i in $(jobs -p); do wait "$i"; done
