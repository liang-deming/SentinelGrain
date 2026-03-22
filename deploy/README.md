# SentinelGrain 部署与联调（工作包 P4）

## 端口与探针（须与配置一致）

| 服务 | 端口 | 说明 |
|------|------|------|
| sentinel-rpc gRPC | 8080 | `ListenOn`（`rpc/etc/sentinel.yaml`） |
| Prometheus metrics | 4000 | `Telemetry.Endpoint` |
| admin HTTP | 8888 | `api/etc/admin.yaml` |
| Redis | 6379 | Admin 与 RPC 须同一实例 |

Kubernetes `livenessProbe` / `readinessProbe` 使用 **TCP 8080**，与 gRPC 监听一致（`开发规范.md` §4.3）。

## 生成 Protobuf（`开发规范.md` §2）

在仓库根目录：

```bash
protoc --proto_path=. --go_out=. --go_opt=module=SentinelGrain --go-grpc_out=. --go-grpc_opt=module=SentinelGrain rpc/sentinel.proto
```

## 本地启动（依赖 Redis）

1. 启动 Redis（默认 `127.0.0.1:6379`）。
2. RPC：`go run ./rpc/cmd/sentinel-rpc -f rpc/etc/sentinel.yaml`
3. Admin：`go run ./api/admin.go -f api/etc/admin.yaml`

## 联调：grpcurl（与 `rpc/sentinel.proto` 一致）

需 RPC 以 **dev** 模式开启 **gRPC reflection**。

```bash
grpcurl -plaintext -d '{"appId":"myapp","resource":"api","dimension":"u1","cost":1}' localhost:8080 sentinel.Sentinel/Check
```

Admin 写入规则：`POST http://localhost:8888/api/v1/quota`（JSON body 见 `api/admin.api`）。**等待 `QuotaRefreshInterval`（默认 5s）** 后再调 Check，规则才会进 RPC 内存表。

## Docker

```bash
docker build -f deploy/docker/Dockerfile -t sentinelgrain/sentinel-rpc:latest .
docker compose -f deploy/docker/docker-compose.yml up --build
```

环境变量：`SENTINEL_REDIS_HOST`（默认 `redis:6379`，适配 compose 内服务名）。

## Kubernetes

```bash
docker build -f deploy/docker/Dockerfile -t sentinelgrain/sentinel-rpc:latest .
kubectl apply -f deploy/k8s/redis.yaml
kubectl apply -f deploy/k8s/sentinel-configmap.yaml
kubectl apply -f deploy/k8s/sentinel-deployment.yaml
kubectl apply -f deploy/k8s/sentinel-service.yaml
```

Deployment 使用 `command` 覆盖镜像内 entrypoint，直接读取 ConfigMap 挂载的 `sentinel.yaml`（Redis 指向集群内 `redis:6379` Service）。

## 测试

```bash
go test ./...
```
