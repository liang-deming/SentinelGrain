# SentinelGrain

[English](README.md) | [日本語](README.ja.md) | [中文](README.zh.md)

SentinelGrain 是一个基于 Go 与 go-zero 实现的轻量级配额/限流服务。
它提供以下能力：

- 用于配额判定的 RPC 服务（`sentinel-rpc`）
- 用于创建和查询配额规则的 Admin HTTP API
- 基于 Redis 的规则存储，并定时刷新到内存缓存
- Docker 与 Kubernetes 部署资产

## 仓库结构

- `rpc/`：gRPC 服务、配额判定逻辑、配置与测试
- `api/`：配额规则管理的 Admin REST 服务
- `deploy/`：Docker/Kubernetes 部署文件与运维说明
- `common/`：公共模型定义与工具代码

## 前置要求

- Go `1.24+`
- Redis（本地开发默认 `127.0.0.1:6379`）
- 可选：用于重新生成 protobuf 的 `protoc`

## 本地运行

1. 启动 Redis。
2. 启动 RPC 服务：

```bash
go run ./rpc/cmd/sentinel-rpc -f rpc/etc/sentinel.yaml
```

3. 启动 Admin 服务：

```bash
go run ./api/admin.go -f api/etc/admin.yaml
```

## API 与 RPC

- Admin API 基础地址：`http://localhost:8888/api/v1`
  - `POST /quota`：创建或更新规则
  - `GET /quota`：查询规则列表
- RPC 地址：`localhost:8080`
  - 服务：`sentinel.Sentinel`
  - 方法：`Check`

`grpcurl` 示例：

```bash
grpcurl -plaintext \
  -d '{"appId":"myapp","resource":"api","dimension":"u1","cost":1}' \
  localhost:8080 sentinel.Sentinel/Check
```

## 配置说明

- RPC 服务器配置位于顶层 `rpc:` 下。
- Redis 业务配置使用顶层 `BizRedis:`。
- `QuotaRefreshInterval` 控制从 Redis 定时刷新规则的周期。
- `L1Cache` 用于热点路径的本地内存缓存。

当前参考配置见 `rpc/etc/sentinel.yaml`。

## 部署

Docker/Kubernetes 详细部署与排障文档见：

- `deploy/README.md`

在本地 Docker Desktop + Kubernetes 环境中，建议使用以下脚本在发布时强制唯一镜像 tag：

```bash
./deploy/k8s/deploy-local-image.sh
```

## 测试

```bash
go test ./...
```
