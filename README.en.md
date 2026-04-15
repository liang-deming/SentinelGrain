# SentinelGrain

SentinelGrain is a lightweight quota/rate-control service built with Go and go-zero.
It provides:

- An RPC service (`sentinel-rpc`) for quota checks.
- An Admin HTTP API for creating and querying quota rules.
- Redis-backed rule storage with periodic refresh into in-memory cache.
- Docker and Kubernetes deployment assets.

## Repository Layout

- `rpc/`: gRPC service, quota check logic, config, tests.
- `api/`: Admin REST service for quota rule management.
- `deploy/`: Docker/Kubernetes deployment files and operational notes.
- `common/`: shared model definitions and helpers.

## Prerequisites

- Go `1.24+`
- Redis (`127.0.0.1:6379` by default for local development)
- Optional: `protoc` for regenerating protobuf files

## Run Locally

1. Start Redis.
2. Start RPC service:

```bash
go run ./rpc/cmd/sentinel-rpc -f rpc/etc/sentinel.yaml
```

3. Start Admin service:

```bash
go run ./api/admin.go -f api/etc/admin.yaml
```

## API and RPC

- Admin API base: `http://localhost:8888/api/v1`
  - `POST /quota` to create/update a rule
  - `GET /quota` to list rules
- RPC endpoint: `localhost:8080`
  - Service: `sentinel.Sentinel`
  - Method: `Check`

Example `grpcurl`:

```bash
grpcurl -plaintext \
  -d '{"appId":"myapp","resource":"api","dimension":"u1","cost":1}' \
  localhost:8080 sentinel.Sentinel/Check
```

## Configuration Notes

- RPC server config is nested under top-level `rpc:`.
- Redis business config uses top-level `BizRedis:`.
- `QuotaRefreshInterval` controls periodic rule refresh from Redis.
- `L1Cache` enables local in-memory cache for hot paths.

See `rpc/etc/sentinel.yaml` for the current reference configuration.

## Deployment

For detailed Docker/Kubernetes setup and troubleshooting, see:

- `deploy/README.md`

Use the helper script below on local Docker Desktop + Kubernetes to force unique image tags during rollout:

```bash
./deploy/k8s/deploy-local-image.sh
```

## Tests

```bash
go test ./...
```
