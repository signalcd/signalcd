#!/usr/bin/env bash

set -euo pipefail

trap 'kill $(jobs -p); exit 0' EXIT

(
  ./development/caddy -conf ./development/Caddyfile
) &

(
  ./cmd/api/api
) &

(
  cd ui && webdev serve web:6670 --tls-cert-chain=../development/signalcd.dev+6.pem --tls-cert-key=../development/signalcd.dev+6-key.pem --auto=restart
) &

for i in `jobs -p`; do wait $i; done
