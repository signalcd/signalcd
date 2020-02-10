#!/usr/bin/env bash

set -euo pipefail

trap 'kill $(jobs -p); exit 0' EXIT

(
  ./development/caddy -conf ./development/Caddyfile-dart
) &

(
  ./cmd/api/api \
      --tls.cert ./development/signalcd.dev+6.pem \
      --tls.key ./development/signalcd.dev+6-key.pem \
      --ui.assets ./ui/build
) &

(
  cd ui && webdev serve web:6670 \
      --tls-cert-chain=../development/signalcd.dev+6.pem \
      --tls-cert-key=../development/signalcd.dev+6-key.pem \
      --auto=restart
) &

for i in `jobs -p`; do wait $i; done
