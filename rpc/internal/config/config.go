package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	// 须具名，勿嵌入：嵌入 zrpc.RpcServerConf 时与 BizRedis 会在 go-zero conf 中触发 conflict key redis。
	// 配置须使用顶层键 `rpc:` 嵌套（勿依赖 inherit 扁平填充，conf.Load 与 zrpc.MustNewServer 行为不一致）。
	Rpc zrpc.RpcServerConf `json:"rpc"`
	BizRedis              redis.RedisConf
	CommandTimeout        int // Redis命令超时时间(毫秒)
	L1Cache               L1Config
	// QuotaRefreshInterval 秒；定时全量从 Redis 加载规则到 RuleCache（P3）。0 时 main 使用默认 5。
	QuotaRefreshInterval int `json:",optional"`
}

type L1Config struct {
	Enabled bool  // 是否启用本地缓存
	TTL     int64 // 缓存TTL(毫秒)
}
