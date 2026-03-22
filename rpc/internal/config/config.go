package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Redis          redis.RedisConf // Redis配置
	CommandTimeout int             // Redis命令超时时间(毫秒)
	L1Cache       L1Config        // 本地缓存配置
	Prometheus    bool            // 是否启用Prometheus指标
}

type L1Config struct {
	Enabled bool  // 是否启用本地缓存
	TTL     int64 // 缓存TTL(毫秒)
}
