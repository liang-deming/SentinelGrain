package errors

// 限流相关错误码
const (
	RateLimitExceeded = "RateLimitExceeded" // 超过限流阈值
	InternalError     = "InternalError"      // 内部错误
	RuleNotFound      = "RuleNotFound"       // 限流规则未找到
)
