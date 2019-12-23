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
  cd ui && webdev serve --auto=restart
) &

for i in `jobs -p`; do wait $i; done
