# SentinelGrain 部署与联调（工作包 P4）

## 端口与探针（须与配置一致）


| 服务                 | 端口   | 说明                                  |
| ------------------ | ---- | ----------------------------------- |
| sentinel-rpc gRPC  | 8080 | `ListenOn`（`rpc/etc/sentinel.yaml`） |
| Prometheus metrics | 4000 | `rpc.Telemetry`（见 `rpc/etc/sentinel.yaml`） |
| admin HTTP         | 8888 | `api/etc/admin.yaml`                |
| Redis              | 6379 | Admin 与 RPC 须同一实例                   |


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

### Docker Desktop：本地镜像「建了但 Pod 还是旧代码」

现象：`kubectl describe pod` 里 **Image ID** 的 digest 在多次 `docker build -t ...:latest` 后不变；事件里常有 **image already present on machine**。此时进程仍是旧二进制（例如一直报 `conflict key redis`）。

**不要**只依赖 `:latest` + `kubectl rollout restart`。请用仓库脚本（每次生成**唯一 tag** 并 `kubectl set image`）：

```bash
./deploy/k8s/deploy-local-image.sh
```

或手动：`docker build ... -t sentinelgrain/sentinel-rpc:sentinel-$(date +%s) .`，再 `kubectl set image deployment/sentinel-rpc sentinel-rpc=<同上镜像名>`。

本地重建镜像后，若 `kubectl apply` 对 Deployment 显示 **unchanged**，集群**不会**自动用新镜像起 Pod，仍在跑旧容器。须让 Deployment 滚动一次，例如：

```bash
kubectl rollout restart deployment/sentinel-rpc
kubectl rollout status deployment/sentinel-rpc --timeout=120s
```

或删除 Pod 由控制器重建：`kubectl delete pod -l app=sentinel-rpc`。

修改 `deploy/k8s/sentinel-configmap.yaml` 后须重新 `kubectl apply -f deploy/k8s/sentinel-configmap.yaml` 并 `kubectl rollout restart deployment/sentinel-rpc`，Pod 才会挂载新配置。

Deployment 使用 `command` 覆盖镜像内 entrypoint，直接读取 ConfigMap 挂载的 `sentinel.yaml`（`rpc` 与 `BizRedis` 同级；Redis 指向集群内 `redis:6379` Service）。

## 测试

```bash
go test ./...
```

