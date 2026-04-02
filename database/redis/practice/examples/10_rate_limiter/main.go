package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	fixedWindowLimiter(ctx, rdb)
	slidingWindowLimiter(ctx, rdb)
	tokenBucketLimiter(ctx, rdb)
	concurrentLimiter(ctx, rdb)
}

// ============================================================
// 1. 固定窗口限流
// ============================================================

var fixedWindowScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

local current = redis.call('INCR', key)
if current == 1 then
    redis.call('EXPIRE', key, window)
end

if current > limit then
    return 0
end
return 1
`)

func fixedWindowLimiter(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 固定窗口限流 ==========")

	key := "ratelimit:fixed:api:/user/info"

	// 限制：每 10 秒最多 5 次
	for i := 1; i <= 8; i++ {
		allowed, _ := fixedWindowScript.Run(ctx, rdb, []string{key}, 5, 10).Int()
		if allowed == 1 {
			fmt.Printf("  请求 #%d: 通过\n", i)
		} else {
			fmt.Printf("  请求 #%d: 被限流\n", i)
		}
	}

	rdb.Del(ctx, key)
	fmt.Println()
}

// ============================================================
// 2. 滑动窗口限流（ZSet 实现）
// ============================================================

var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local member = ARGV[4]

-- 移除窗口外的过期记录
redis.call('ZREMRANGEBYSCORE', key, '-inf', now - window)

-- 统计窗口内的请求数
local count = redis.call('ZCARD', key)

if count >= limit then
    return 0
end

-- 添加当前请求
redis.call('ZADD', key, now, member)
redis.call('PEXPIRE', key, window)

return 1
`)

func slidingWindowLimiter(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 滑动窗口限流 ==========")

	key := "ratelimit:sliding:api:/order/create"

	// 限制：1000ms 窗口内最多 3 次
	window := 1000 // 毫秒
	limit := 3

	for i := 1; i <= 6; i++ {
		now := time.Now().UnixMilli()
		member := fmt.Sprintf("%d-%d", now, i) // 唯一成员

		allowed, _ := slidingWindowScript.Run(ctx, rdb,
			[]string{key}, limit, window, now, member).Int()
		if allowed == 1 {
			fmt.Printf("  请求 #%d: 通过 (t=%dms)\n", i, now%10000)
		} else {
			fmt.Printf("  请求 #%d: 被限流 (t=%dms)\n", i, now%10000)
		}
		time.Sleep(200 * time.Millisecond)
	}

	rdb.Del(ctx, key)
	fmt.Println()
}

// ============================================================
// 3. 令牌桶限流（Lua 实现）
// ============================================================

var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]
local rate = tonumber(ARGV[1])      -- 每秒生成令牌数
local capacity = tonumber(ARGV[2])  -- 桶容量
local now = tonumber(ARGV[3])       -- 当前时间（毫秒）
local requested = tonumber(ARGV[4]) -- 请求消耗的令牌数

-- 获取上次状态
local last_time = tonumber(redis.call('HGET', key, 'last_time') or now)
local tokens = tonumber(redis.call('HGET', key, 'tokens') or capacity)

-- 计算新增令牌
local elapsed = (now - last_time) / 1000
local new_tokens = elapsed * rate
tokens = math.min(capacity, tokens + new_tokens)

-- 判断是否放行
if tokens >= requested then
    tokens = tokens - requested
    redis.call('HSET', key, 'last_time', now, 'tokens', tokens)
    redis.call('PEXPIRE', key, 60000)
    return 1
end

redis.call('HSET', key, 'last_time', now, 'tokens', tokens)
redis.call('PEXPIRE', key, 60000)
return 0
`)

func tokenBucketLimiter(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 令牌桶限流 ==========")

	key := "ratelimit:token:api:/pay"

	// 配置：每秒 2 个令牌，桶容量 5
	rate := 2
	capacity := 5

	// 瞬间发 6 个请求（桶初始满 = 5 个令牌）
	fmt.Println("--- 瞬间 burst ---")
	for i := 1; i <= 6; i++ {
		now := time.Now().UnixMilli()
		allowed, _ := tokenBucketScript.Run(ctx, rdb,
			[]string{key}, rate, capacity, now, 1).Int()
		if allowed == 1 {
			fmt.Printf("  请求 #%d: 通过\n", i)
		} else {
			fmt.Printf("  请求 #%d: 被限流（令牌不足）\n", i)
		}
	}

	// 等 1 秒，生成 2 个令牌
	fmt.Println("--- 等 1 秒后 ---")
	time.Sleep(1 * time.Second)
	for i := 7; i <= 10; i++ {
		now := time.Now().UnixMilli()
		allowed, _ := tokenBucketScript.Run(ctx, rdb,
			[]string{key}, rate, capacity, now, 1).Int()
		if allowed == 1 {
			fmt.Printf("  请求 #%d: 通过\n", i)
		} else {
			fmt.Printf("  请求 #%d: 被限流\n", i)
		}
	}

	rdb.Del(ctx, key)
	fmt.Println()
}

// ============================================================
// 4. 并发压测限流器
// ============================================================
func concurrentLimiter(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 并发压测限流器 ==========")

	key := "ratelimit:bench:api"

	// 50 个并发请求，窗口 1s，限制 20 QPS
	var (
		wg      sync.WaitGroup
		passed  int64
		blocked int64
	)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			now := time.Now().UnixMilli()
			member := fmt.Sprintf("%d-%d", now, id)
			allowed, _ := slidingWindowScript.Run(ctx, rdb,
				[]string{key}, 20, 1000, now, member).Int()
			if allowed == 1 {
				atomic.AddInt64(&passed, 1)
			} else {
				atomic.AddInt64(&blocked, 1)
			}
		}(i)
	}
	wg.Wait()

	fmt.Printf("  50 并发请求，限制 20 QPS:\n")
	fmt.Printf("  通过: %d, 被限流: %d\n", passed, blocked)

	rdb.Del(ctx, key)
	fmt.Println()

	log.Println("[10_rate_limiter] 限流器演示完成")
}
