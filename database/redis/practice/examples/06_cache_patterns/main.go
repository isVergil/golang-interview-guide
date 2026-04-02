package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

// 演示三种缓存问题的解决方案：穿透、击穿、旁路缓存
func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	cacheAside(ctx, rdb)
	cacheNullValue(ctx, rdb)
	cacheBreakdownMutex(ctx, rdb)
	cacheBreakdownLogicalExpire(ctx, rdb)
}

// 模拟数据库查询
func queryDB(id int) (string, bool) {
	// 模拟耗时
	time.Sleep(50 * time.Millisecond)
	db := map[int]string{
		1: `{"name":"iPhone","price":5999}`,
		2: `{"name":"MacBook","price":12999}`,
	}
	data, ok := db[id]
	return data, ok
}

// ============================================================
// 1. Cache Aside（旁路缓存）：标准读写模式
// ============================================================
func cacheAside(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Cache Aside 旁路缓存 ==========")

	key := "cache:product:1"

	// --- 读路径（懒加载） ---
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// 缓存 miss → 查 DB → 回填缓存
		fmt.Println("缓存 miss，查 DB...")
		data, ok := queryDB(1)
		if ok {
			// 过期时间加随机抖动，防雪崩
			ttl := 3600*time.Second + time.Duration(time.Now().UnixNano()%300)*time.Second
			rdb.Set(ctx, key, data, ttl)
			fmt.Println("回填缓存, data =", data)
		}
	} else {
		fmt.Println("缓存 hit, data =", val)
	}

	// 第二次读：命中缓存
	val, _ = rdb.Get(ctx, key).Result()
	fmt.Println("第二次读（缓存 hit） =", val)

	// --- 写路径 ---
	// 先更新 DB（这里省略），再删除缓存
	fmt.Println("更新 DB 后，删除缓存")
	rdb.Del(ctx, key)

	// 下次读自动触发懒加载
	val, err = rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("缓存已删，下次读触发懒加载")
	}

	rdb.Del(ctx, key)
	fmt.Println()
}

// ============================================================
// 2. 缓存穿透：缓存空值方案
// ============================================================
func cacheNullValue(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 缓存穿透：缓存空值 ==========")

	key := "cache:product:9999" // 不存在的 ID

	for i := 1; i <= 3; i++ {
		val, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			// 缓存 miss → 查 DB
			data, ok := queryDB(9999)
			if !ok {
				// DB 也没有 → 缓存空值，设短过期（防止缓存污染）
				rdb.Set(ctx, key, "", 60*time.Second)
				fmt.Printf("  第 %d 次：DB 查不到，缓存空值（TTL=60s）\n", i)
			} else {
				rdb.Set(ctx, key, data, 3600*time.Second)
			}
		} else if val == "" {
			// 命中空值缓存 → 直接返回，不查 DB
			fmt.Printf("  第 %d 次：命中空值缓存，不查 DB\n", i)
		} else {
			fmt.Printf("  第 %d 次：缓存 hit = %s\n", i, val)
		}
	}

	rdb.Del(ctx, key)
	fmt.Println()
}

// ============================================================
// 3. 缓存击穿：互斥锁方案
// ============================================================
func cacheBreakdownMutex(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 缓存击穿：互斥锁 ==========")

	key := "cache:hot:1"
	lockKey := "lock:cache:hot:1"

	// 模拟 5 个并发请求同时发现缓存失效
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			val, err := rdb.Get(ctx, key).Result()
			if err == redis.Nil {
				// 缓存 miss → 抢锁
				ok, _ := rdb.SetNX(ctx, lockKey, 1, 10*time.Second).Result()
				if ok {
					// 拿到锁 → 查 DB → 回填缓存 → 释放锁
					fmt.Printf("  请求 #%d: 拿到锁，查 DB\n", id)
					data, _ := queryDB(1)
					rdb.Set(ctx, key, data, time.Minute)
					rdb.Del(ctx, lockKey)
					fmt.Printf("  请求 #%d: 回填缓存完成\n", id)
				} else {
					// 没拿到锁 → 等待重试
					fmt.Printf("  请求 #%d: 等待中...\n", id)
					time.Sleep(100 * time.Millisecond)
					val, _ = rdb.Get(ctx, key).Result()
					fmt.Printf("  请求 #%d: 重试获取到 = %s\n", id, val[:20]+"...")
				}
			} else {
				fmt.Printf("  请求 #%d: 缓存 hit = %s\n", id, val[:20]+"...")
				_ = val
			}
		}(i)
	}
	wg.Wait()

	rdb.Del(ctx, key, lockKey)
	fmt.Println()
}

// ============================================================
// 4. 缓存击穿：逻辑过期方案（推荐）
// ============================================================

type CacheData struct {
	Data     string `json:"data"`
	ExpireAt int64  `json:"expire_at"` // 逻辑过期时间戳
}

func cacheBreakdownLogicalExpire(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 缓存击穿：逻辑过期 ==========")

	key := "cache:logical:1"
	lockKey := "lock:refresh:logical:1"

	// 预热：写入数据，逻辑过期时间设为"已过期"
	data := CacheData{
		Data:     `{"name":"iPhone","price":5999}`,
		ExpireAt: time.Now().Add(-1 * time.Second).Unix(), // 已过期
	}
	bytes, _ := json.Marshal(data)
	rdb.Set(ctx, key, string(bytes), 0) // Redis key 永不过期

	// 模拟 3 个并发请求
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 一定能 hit（key 永不过期）
			val, _ := rdb.Get(ctx, key).Result()
			var cached CacheData
			json.Unmarshal([]byte(val), &cached)

			if cached.ExpireAt < time.Now().Unix() {
				// 逻辑过期了 → 抢锁异步刷新
				ok, _ := rdb.SetNX(ctx, lockKey, 1, 10*time.Second).Result()
				if ok {
					// 拿到锁 → 异步刷新（不阻塞当前请求）
					go func() {
						fmt.Printf("  请求 #%d: 拿到锁，异步刷新缓存...\n", id)
						newData, _ := queryDB(1)
						refreshed := CacheData{
							Data:     newData,
							ExpireAt: time.Now().Add(3600 * time.Second).Unix(),
						}
						bytes, _ := json.Marshal(refreshed)
						rdb.Set(ctx, key, string(bytes), 0)
						rdb.Del(ctx, lockKey)
						fmt.Printf("  请求 #%d: 异步刷新完成\n", id)
					}()
				}
				// 不管有没有拿到锁，都返回旧数据（不阻塞）
				fmt.Printf("  请求 #%d: 返回旧数据 = %s\n", id, cached.Data)
			} else {
				fmt.Printf("  请求 #%d: 数据未过期 = %s\n", id, cached.Data)
			}
		}(i)
	}
	wg.Wait()
	time.Sleep(200 * time.Millisecond) // 等异步刷新完成

	rdb.Del(ctx, key, lockKey)
	fmt.Println()

	log.Println("[06_cache_patterns] 缓存模式演示完成")
}
