package quota

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reAppID    = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)
	reResource = regexp.MustCompile(`^[a-zA-Z0-9_./-]{1,128}$`)
)

// ValidateRuleKey 校验 appId、resource（用于 Redis key 段，禁止 : 与索引分隔符）
func ValidateRuleKey(appId, resource string) error {
	if !reAppID.MatchString(appId) {
		return fmt.Errorf("appId: invalid format (allowed: letters, digits, _, -, length 1-64)")
	}
	if !reResource.MatchString(resource) {
		return fmt.Errorf("resource: invalid format (allowed: letters, digits, _, ., /, -, length 1-128)")
	}
	if strings.Contains(appId, ":") || strings.Contains(resource, ":") {
		return fmt.Errorf("appId and resource must not contain ':'")
	}
	if strings.Contains(appId, indexSep) || strings.Contains(resource, indexSep) {
		return fmt.Errorf("appId and resource must not contain %q", indexSep)
	}
	return nil
}

// ValidateRule 校验阈值与窗口
func ValidateRule(threshold, period int64) error {
	if threshold <= 0 {
		return fmt.Errorf("threshold must be positive")
	}
	if period <= 0 {
		return fmt.Errorf("period (seconds) must be positive")
	}
	return nil
}
