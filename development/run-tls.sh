#!/usr/bin/env bash

set -euo pipefail

trap 'kill $(jobs -p); exit 0' EXIT

(
  ./development/caddy -conf ./development/Caddyfile
) &

(
  ./cmd/api/api --tls.cert ./development/signalcd.dev+6.pem --tls.key ./development/signalcd.dev+6-key.pem
) &

(
  ./cmd/ui/ui
) &

for i in `jobs -p`; do wait $i; done
