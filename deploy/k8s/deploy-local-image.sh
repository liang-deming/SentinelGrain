#!/usr/bin/env bash
# Docker Desktop 内置 Kubernetes：反复打 tag :latest 时，kubelet 常仍使用 containerd 里旧的 digest，
# Pod 里跑旧二进制（日志会一直 conflict key redis）。请用「每次唯一 tag」强制换新镜像。
# 用法（在仓库根目录）: ./deploy/k8s/deploy-local-image.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TAG="sentinel-$(date +%s)"
IMAGE="sentinelgrain/sentinel-rpc:${TAG}"

echo "==> docker build -> ${IMAGE}"
docker build --no-cache -f "${ROOT}/deploy/docker/Dockerfile" -t "${IMAGE}" "${ROOT}"

echo "==> kubectl set image (forces new ReplicaSet + new digest)"
kubectl set image deployment/sentinel-rpc sentinel-rpc="${IMAGE}"

echo "==> wait rollout"
kubectl rollout status deployment/sentinel-rpc --timeout=180s

echo "==> done. Current image:"
kubectl get deploy sentinel-rpc -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'
