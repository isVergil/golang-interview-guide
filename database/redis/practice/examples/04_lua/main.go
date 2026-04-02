package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

// 演示 Lua 脚本在 Redis 中的原子操作
func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	luaCompareAndDelete(ctx, rdb)
	luaRateLimiter(ctx, rdb)
	luaStockDeduct(ctx, rdb)
	luaCached(ctx, rdb)
}

// ============================================================
// 1. CAS 删除：先比较再删除（分布式锁释放的标准做法）
// ============================================================
func luaCompareAndDelete(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Lua: CAS 删除（锁释放） ==========")

	// 释放锁的 Lua 脚本：先判断值是不是自己的，是才删
	// 保证 GET + DEL 的原子性，防止误删别人的锁
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)

	// 模拟加锁
	rdb.Set(ctx, "lock:order:1", "holder-A", 0)

	// holder-B 尝试释放 → 失败（不是自己的锁）
	result, _ := script.Run(ctx, rdb, []string{"lock:order:1"}, "holder-B").Int64()
	fmt.Println("holder-B 释放结果 =", result) // 0

	// holder-A 释放 → 成功
	result, _ = script.Run(ctx, rdb, []string{"lock:order:1"}, "holder-A").Int64()
	fmt.Println("holder-A 释放结果 =", result) // 1

	exists, _ := rdb.Exists(ctx, "lock:order:1").Result()
	fmt.Println("锁是否还在 =", exists) // 0
	fmt.Println()
}

// ============================================================
// 2. 固定窗口限流
// ============================================================
func luaRateLimiter(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Lua: 固定窗口限流 ==========")

	// Lua 脚本：INCR + 设置过期时间，保证原子性
	// KEYS[1]=限流 key, ARGV[1]=最大请求数, ARGV[2]=窗口时间(秒)
	//
	// 注意：只有通过此脚本执行才会计数（内部 INCR）
	// 外部直接 GET 该 key 不影响计数；但外部 INCR/SET 会干扰计数
	// 所以限流 key 不要被其他业务共用
	script := redis.NewScript(`
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local current = tonumber(redis.call("GET", key) or "0")
		if current + 1 > limit then
			return 0
		end
		redis.call("INCR", key)
		if current == 0 then
			redis.call("EXPIRE", key, window)
		end
		return 1
	`)

	key := "rate:user:1:api"
	limit := 5
	window := 60

	for i := 1; i <= 7; i++ {
		allowed, _ := script.Run(ctx, rdb, []string{key}, limit, window).Int64()
		if allowed == 1 {
			fmt.Printf("  请求 #%d → 放行\n", i)
		} else {
			fmt.Printf("  请求 #%d → 限流（超过 %d 次/分钟）\n", i, limit)
		}
	}

	rdb.Del(ctx, key)
	fmt.Println()
}

// ============================================================
// 3. 库存扣减（原子性保证不超卖）
// ============================================================
func luaStockDeduct(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Lua: 库存扣减 ==========")

	// Lua 脚本：检查库存 >= 扣减量才扣，防止超卖
	script := redis.NewScript(`
		local stock = tonumber(redis.call("GET", KEYS[1]) or "0")
		local deduct = tonumber(ARGV[1])
		if stock >= deduct then
			redis.call("DECRBY", KEYS[1], deduct)
			return stock - deduct
		else
			return -1
		end
	`)

	rdb.Set(ctx, "stock:product:1", 10, 0)

	// 扣 3 个
	remain, _ := script.Run(ctx, rdb, []string{"stock:product:1"}, 3).Int64()
	fmt.Println("扣 3 后库存 =", remain) // 7

	// 扣 5 个
	remain, _ = script.Run(ctx, rdb, []string{"stock:product:1"}, 5).Int64()
	fmt.Println("扣 5 后库存 =", remain) // 2

	// 扣 5 个（库存不足）
	remain, _ = script.Run(ctx, rdb, []string{"stock:product:1"}, 5).Int64()
	fmt.Println("扣 5 结果 =", remain) // -1（库存不足）

	rdb.Del(ctx, "stock:product:1")
	fmt.Println()
}

// ============================================================
// 4. EVALSHA 缓存脚本（生产优化）
// ============================================================
func luaCached(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Lua: EVALSHA 脚本缓存 ==========")

	// redis.NewScript 内部自动处理 EVALSHA/EVAL 降级：
	//
	// EVAL  = 每次传完整脚本文本给 Redis 执行（网络开销大）
	// EVALSHA = 只传脚本的 SHA1 哈希值（40 字节），Redis 用哈希查缓存执行
	//
	// 完整流程：
	//   第 1 次 Run()：
	//     客户端 → EVALSHA sha1（只发哈希）
	//     Redis  → NOSCRIPT（没见过这个脚本）
	//     客户端 → EVAL "完整脚本"（降级，发原文）
	//     Redis  → 执行 + 把脚本缓存在内存（key=SHA1, value=脚本原文）
	//
	//   第 2 次 Run()：
	//     客户端 → EVALSHA sha1（只发哈希）
	//     Redis  → 命中缓存，直接执行（不需要再传脚本）
	//
	// 注意：Redis 重启后脚本缓存丢失，下次 EVALSHA 会自动降级回 EVAL 重新缓存
	script := redis.NewScript(`return redis.call("SET", KEYS[1], ARGV[1])`)

	// 第一次：EVALSHA miss → EVAL（传完整脚本）
	script.Run(ctx, rdb, []string{"lua:cached:test"}, "hello")

	// 第二次：EVALSHA hit（只传 SHA1，不传脚本体，省带宽）
	script.Run(ctx, rdb, []string{"lua:cached:test"}, "world")

	val, _ := rdb.Get(ctx, "lua:cached:test").Result()
	fmt.Println("lua:cached:test =", val) // world

	rdb.Del(ctx, "lua:cached:test")
	fmt.Println()

	log.Println("[04_lua] Lua 脚本演示完成")
}
