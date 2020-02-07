# Development

Start the API and Agent locally:
```bash
make cmd/api/api && ./cmd/api/api
make cmd/agent/agent && ./cmd/agent/agent -kubeconfig ~/.kube/some-cluster -name local
```

## UI

If you want to start the UI and develop in Dart, you need to run a proxy in front of the API and UI
to circumvent Cross-Origin Resource Sharing (CORS).

For that reason we use [Caddy](https://caddyserver.com/).
In the Caddyfile you can see a pre-configured server running on 6060 that forwards on `/` to the UI
and everything starting with `/api` is forwarded to the API.

[Download Caddy](https://caddyserver.com/download) and extract the `caddy` binary
into this directory and start it using by simply running `./caddy`.
The provided `Caddyfile` will be used automatically.


## Ports

| Component             | Port  |
|-----------------------|-------|
| Caddy for Dev         | 6060  |
| API HTTP/gRPC Public  | 6660  |
| API gRPC Agent        | 6661  |
| API HTTP Internal     | 6662  |
| UI HTTP Public        | 6670  |
