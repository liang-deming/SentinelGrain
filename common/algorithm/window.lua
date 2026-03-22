-- KEYS[1]: 限流key
-- ARGV[1]: 当前时间戳(ms)
-- ARGV[2]: 窗口大小(ms)
-- ARGV[3]: 阈值
-- ARGV[4]: 本次消耗配额

local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local threshold = tonumber(ARGV[3])
local cost = tonumber(ARGV[4])

-- 1. 移除窗口外的数据
local windowStart = now - window
redis.call('ZREMRANGEBYSCORE', key, '-inf', windowStart)

-- 2. 获取窗口内当前计数
local current = redis.call('ZCARD', key)

-- 3. 判断是否允许通过
local allowed = (current + cost) <= threshold
local remaining = threshold - current

if allowed then
    -- 4. 添加当前请求到集合
    local member = now .. ':' .. redis.call('INCR', key .. ':seq')
    redis.call('ZADD', key, now, member)
    remaining = remaining - cost
end

-- 5. 返回结构化结果
return {
    allowed and 1 or 0,  -- allowed
    remaining            -- remaining
}
