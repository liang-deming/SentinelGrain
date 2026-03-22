package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Redis                 redis.RedisConf // Redis配置
	CommandTimeout        int             // Redis命令超时时间(毫秒)
	L1Cache               L1Config        // 本地缓存配置
	Prometheus            bool            // 是否启用Prometheus指标
	QuotaRefreshInterval  int             `json:",optional"` // 秒，规则表从 Redis 全量刷新的周期；0 表示使用默认 5
}

type L1Config struct {
	Enabled bool  // 是否启用本地缓存
	TTL     int64 // 缓存TTL(毫秒)
}
