# 多阶段构建 sentinel-rpc（工作包 P4 / 开发规范 4.x）
# 构建：在仓库根目录执行
#   docker build -f deploy/docker/Dockerfile -t sentinelgrain/sentinel-rpc:latest .

FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/sentinel-rpc ./rpc/cmd/sentinel-rpc

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata \
  && addgroup -g 65532 nonroot \
  && adduser -D -u 65532 -G nonroot nonroot \
  && mkdir -p /app/etc
WORKDIR /app
COPY deploy/docker/entrypoint.sh /app/entrypoint.sh
COPY deploy/docker/sentinel.container.yaml /app/etc/sentinel.yaml.template
COPY --from=builder /out/sentinel-rpc /app/sentinel-rpc
RUN chmod +x /app/entrypoint.sh /app/sentinel-rpc \
  && chown -R nonroot:nonroot /app
USER nonroot:nonroot
EXPOSE 8080 4000
ENV SENTINEL_REDIS_HOST=redis:6379
ENTRYPOINT ["/app/entrypoint.sh"]
