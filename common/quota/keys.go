// Package quota 定义 Admin API 与 RPC Check 共用的配额规则存储约定（MVP：Redis）。
package quota

import "fmt"

// Redis Key 约定（单一真相，Admin 与 RPC 共用）：
// - 每条规则一个 Hash：KeyRulePrefix + appId + ":" + resource（appId/resource 经 ValidateRuleKey 校验，不含冒号）
// - 索引 SET：KeyIndex，member 为 IndexMember(appId, resource)，用于 List（MVP 小数据量全量索引）
// - 版本号 STRING：KeyVersion，INCR 用于热更；RPC RuleCache 读此版本并拼入 L1 key，避免 Admin 改阈值后长期脏读

const (
	KeyPrefix = "sentinel:quota:v1"

	fieldThreshold  = "threshold"
	fieldPeriod       = "period"
	fieldUpdatedAt    = "updated_at"
	indexSep          = "|" // appId 与 resource 中禁止出现该字符
)

// KeyRule 单条规则 Redis Hash 的 key：sentinel:quota:v1:rule:{appId}:{resource}
func KeyRule(appId, resource string) string {
	return fmt.Sprintf("%s:rule:%s:%s", KeyPrefix, appId, resource)
}

// KeyIndex 全部规则索引（SET），member 格式见 IndexMember
func KeyIndex() string {
	return KeyPrefix + ":idx"
}

// KeyVersion 全局配置版本，Save 时 INCR；RPC 定时 Refresh 并读此值参与 L1 cache key
func KeyVersion() string {
	return KeyPrefix + ":ver"
}

// IndexMember 索引 SET 中的 member，与 KeyRule 可互推（split 后校验）
func IndexMember(appId, resource string) string {
	return appId + indexSep + resource
}
