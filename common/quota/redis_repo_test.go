package quota

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestRedisRepo_SaveList(t *testing.T) {
	s := miniredis.RunT(t)
	defer s.Close()

	r := NewRedisRepo(redis.New(s.Addr()))
	ctx := context.Background()

	assert.NoError(t, r.Save(ctx, &Rule{AppId: "app1", Resource: "r1", Threshold: 10, Period: 2}))
	assert.NoError(t, r.Save(ctx, &Rule{AppId: "app2", Resource: "r2", Threshold: 3, Period: 1}))

	rules, total, err := r.List(ctx, ListQuery{Page: 1, Size: 10})
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, rules, 2)

	rules, total, err = r.List(ctx, ListQuery{AppId: "app1", Page: 1, Size: 10})
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, "app1", rules[0].AppId)
}
