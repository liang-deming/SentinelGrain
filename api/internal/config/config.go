// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Redis               redis.RedisConf
	RedisCommandTimeout int `json:",default=3000"` // 毫秒，Admin Redis 命令超时
}
