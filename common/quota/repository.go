package quota

import "context"

// Repository 配额持久化抽象，MVP 为 Redis；后续可换 DB 实现
type Repository interface {
	Save(ctx context.Context, rule *Rule) error
	List(ctx context.Context, q ListQuery) ([]Rule, int, error)
}
