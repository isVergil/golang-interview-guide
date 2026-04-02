package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

// 演示生产级分布式锁：SETNX 加锁 + Lua 释放 + 看门狗续期
func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	basicLock(ctx, rdb)
	lockWithWatchdog(ctx, rdb)
	concurrentLock(ctx, rdb)
}

// ============================================================
// 释放锁的 Lua 脚本（全局复用）
// ============================================================
var unlockScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

// ============================================================
// 1. 基础分布式锁：SETNX 加锁 + Lua 释放
// ============================================================
func basicLock(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 基础分布式锁 ==========")

	lockKey := "lock:order:1"
	lockValue := uuid.New().String() // 用 UUID 标识持锁者
	lockTTL := 30 * time.Second

	// 加锁：SET key value NX EX（原子操作）
	ok, err := rdb.SetNX(ctx, lockKey, lockValue, lockTTL).Result()
	if err != nil {
		log.Fatalf("加锁失败: %v", err)
	}
	if !ok {
		fmt.Println("加锁失败：锁已被持有")
		return
	}
	fmt.Println("加锁成功, value =", lockValue)

	// 模拟业务处理
	fmt.Println("处理业务中...")
	time.Sleep(100 * time.Millisecond)

	// 释放锁：Lua 脚本保证原子性（先判断是不是自己的锁再删）
	result, _ := unlockScript.Run(ctx, rdb, []string{lockKey}, lockValue).Int64()
	if result == 1 {
		fmt.Println("释放锁成功")
	} else {
		fmt.Println("释放锁失败：锁已过期或被其他人持有")
	}
	fmt.Println()
}

// ============================================================
// 2. 看门狗续期：防止锁过期但业务没执行完
// ============================================================
func lockWithWatchdog(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 看门狗续期 ==========")

	lockKey := "lock:watchdog:1"
	lockValue := uuid.New().String()
	lockTTL := 5 * time.Second // 初始 5 秒，靠看门狗续期

	// 加锁
	ok, _ := rdb.SetNX(ctx, lockKey, lockValue, lockTTL).Result()
	if !ok {
		fmt.Println("加锁失败")
		return
	}
	fmt.Println("加锁成功, TTL =", lockTTL)

	// 启动看门狗 goroutine：每隔 TTL/3 续期一次
	watchdogCtx, cancelWatchdog := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(lockTTL / 3) // 每 ~1.6 秒续期
		defer ticker.Stop()
		for {
			select {
			case <-watchdogCtx.Done():
				fmt.Println("  看门狗停止")
				return
			case <-ticker.C:
				// 续期前先检查锁还是不是自己的
				val, err := rdb.Get(ctx, lockKey).Result()
				if err != nil || val != lockValue {
					fmt.Println("  看门狗：锁已丢失，停止续期")
					return
				}
				rdb.Expire(ctx, lockKey, lockTTL)
				newTTL, _ := rdb.TTL(ctx, lockKey).Result()
				fmt.Printf("  看门狗续期: TTL 重置为 %v\n", newTTL)
			}
		}
	}()

	// 模拟耗时业务（超过初始 TTL）
	fmt.Println("处理耗时业务（3秒，看门狗会续期）...")
	time.Sleep(3 * time.Second)

	// 业务完成，停止看门狗，释放锁
	cancelWatchdog()
	wg.Wait()

	result, _ := unlockScript.Run(ctx, rdb, []string{lockKey}, lockValue).Int64()
	fmt.Println("释放锁结果 =", result) // 1
	fmt.Println()
}

// ============================================================
// 3. 并发抢锁演示
// ============================================================
func concurrentLock(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 并发抢锁 ==========")

	lockKey := "lock:concurrent:1"
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			lockValue := fmt.Sprintf("worker-%d", id)

			// 尝试加锁
			ok, _ := rdb.SetNX(ctx, lockKey, lockValue, 5*time.Second).Result()
			if !ok {
				fmt.Printf("  worker-%d: 抢锁失败\n", id)
				return
			}
			fmt.Printf("  worker-%d: 抢锁成功，处理业务中...\n", id)

			time.Sleep(200 * time.Millisecond)

			// 释放锁
			unlockScript.Run(ctx, rdb, []string{lockKey}, lockValue)
			fmt.Printf("  worker-%d: 释放锁完成\n", id)
		}(i)
	}

	wg.Wait()
	rdb.Del(ctx, lockKey)
	fmt.Println()

	log.Println("[05_distributed_lock] 分布式锁演示完成")
}
