package config

import (
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

// 与 deploy/k8s/sentinel-configmap.yaml 一致（顶层 rpc: 嵌套）。
func TestLoadYamlLikeK8sConfigMap(t *testing.T) {
	yaml := `
rpc:
  Name: sentinel.rpc
  ListenOn: 0.0.0.0:8080
  Mode: dev
  Telemetry:
    Name: sentinel.rpc
    Endpoint: 0.0.0.0:4000
    Path: /metrics
BizRedis:
  Host: redis:6379
  Type: node
  Pass: ""
  Tls: false
CommandTimeout: 100
QuotaRefreshInterval: 5
L1Cache:
  Enabled: true
  TTL: 500
`
	var c Config
	if err := conf.LoadFromYamlBytes([]byte(yaml), &c); err != nil {
		t.Fatal(err)
	}
	if c.Rpc.ListenOn != "0.0.0.0:8080" {
		t.Fatalf("Rpc.ListenOn: got %q", c.Rpc.ListenOn)
	}
	if c.BizRedis.Host != "redis:6379" {
		t.Fatalf("BizRedis.Host: got %q", c.BizRedis.Host)
	}
}
