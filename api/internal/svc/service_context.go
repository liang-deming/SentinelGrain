// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"SentinelGrain/api/internal/config"
	"SentinelGrain/common/quota"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config    config.Config
	BizRedis  *redis.Redis
	QuotaRepo quota.Repository
}

func NewServiceContext(c config.Config) *ServiceContext {
	r := redis.MustNewRedis(c.Redis)
	return &ServiceContext{
		Config:    c,
		BizRedis:  r,
		QuotaRepo: quota.NewRedisRepo(r),
	}
}
