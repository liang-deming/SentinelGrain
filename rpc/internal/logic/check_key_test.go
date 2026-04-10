package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimiterKeyFormat(t *testing.T) {
	// 与 Admin 写入的 appId/resource 语义一致；dimension 为 Check 维度；v 后为 RuleCache 读入的配额版本
	assert.Equal(t, "myapp:/api/v1:user-1:v7", LimiterKey("myapp", "/api/v1", "user-1", 7))
	assert.Equal(t, "a:b:c:v0", LimiterKey("a", "b", "c", 0))
}
