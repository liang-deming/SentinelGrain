package quota

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RedisRepo MVP 实现 Repository，直连 Redis（方案 B）；后续可替换为 DB 实现同一接口
type RedisRepo struct {
	r *redis.Redis
}

// NewRedisRepo 由 svc 注入的客户端构造，禁止在 Logic 中直接 NewRedis
func NewRedisRepo(r *redis.Redis) *RedisRepo {
	return &RedisRepo{r: r}
}

// Save 写入规则 Hash、索引 SET，并 INCR 全局版本号以触发 RPC 侧热更与 L1 key 隔离
func (r *RedisRepo) Save(ctx context.Context, rule *Rule) error {
	if rule == nil {
		return ErrNilRule
	}
	if err := ValidateRuleKey(rule.AppId, rule.Resource); err != nil {
		return err
	}
	if err := ValidateRule(rule.Threshold, rule.Period); err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	rule.UpdateTime = now
	ruleKey := KeyRule(rule.AppId, rule.Resource)
	fields := map[string]string{
		fieldThreshold:  strconv.FormatInt(rule.Threshold, 10),
		fieldPeriod:     strconv.FormatInt(rule.Period, 10),
		fieldUpdatedAt: strconv.FormatInt(now, 10),
	}
	if err := r.r.HmsetCtx(ctx, ruleKey, fields); err != nil {
		return err
	}
	if _, err := r.r.SaddCtx(ctx, KeyIndex(), IndexMember(rule.AppId, rule.Resource)); err != nil {
		return err
	}
	if _, err := r.r.IncrCtx(ctx, KeyVersion()); err != nil {
		return err
	}
	return nil
}

// List 按索引列出规则，支持可选过滤与分页（MVP：全量索引内存过滤，适合小数据量）
func (r *RedisRepo) List(ctx context.Context, q ListQuery) ([]Rule, int, error) {
	members, err := r.r.SmembersCtx(ctx, KeyIndex())
	if err != nil {
		return nil, 0, err
	}
	var rules []Rule
	for _, m := range members {
		parts := strings.SplitN(m, indexSep, 2)
		if len(parts) != 2 {
			continue
		}
		appId, res := parts[0], parts[1]
		if q.AppId != "" && appId != q.AppId {
			continue
		}
		if q.Resource != "" && res != q.Resource {
			continue
		}
		rk := KeyRule(appId, res)
		hm, err := r.r.HgetallCtx(ctx, rk)
		if err != nil {
			return nil, 0, err
		}
		rule, ok := ruleFromHash(appId, res, hm)
		if !ok {
			continue
		}
		rules = append(rules, rule)
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].UpdateTime > rules[j].UpdateTime })
	total := len(rules)
	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.Size
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	start := int((page - 1) * size)
	if start >= total {
		return []Rule{}, total, nil
	}
	end := start + int(size)
	if end > total {
		end = total
	}
	return rules[start:end], total, nil
}

func ruleFromHash(appId, resource string, hm map[string]string) (Rule, bool) {
	if len(hm) == 0 {
		return Rule{}, false
	}
	th, err1 := strconv.ParseInt(hm[fieldThreshold], 10, 64)
	pd, err2 := strconv.ParseInt(hm[fieldPeriod], 10, 64)
	ut, err3 := strconv.ParseInt(hm[fieldUpdatedAt], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return Rule{}, false
	}
	return Rule{
		AppId:      appId,
		Resource:   resource,
		Threshold:  th,
		Period:     pd,
		UpdateTime: ut,
	}, true
}
