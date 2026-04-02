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

// 基于 ZSet 实现延迟队列：Score = 执行时间戳
func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	basicDelayedQueue(ctx, rdb)
	producerConsumer(ctx, rdb)
}

// ============================================================
// 1. 基础延迟队列
// ============================================================

type DelayedTask struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func basicDelayedQueue(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 基础延迟队列 ==========")

	queueKey := "delayed:queue:basic"

	// 生产：添加延迟任务，Score = 执行时间戳
	tasks := []struct {
		task  DelayedTask
		delay time.Duration
	}{
		{DelayedTask{ID: "t1", Type: "email", Payload: "发送欢迎邮件"}, 1 * time.Second},
		{DelayedTask{ID: "t2", Type: "sms", Payload: "发送验证码"}, 500 * time.Millisecond},
		{DelayedTask{ID: "t3", Type: "callback", Payload: "回调通知"}, 2 * time.Second},
	}

	for _, t := range tasks {
		data, _ := json.Marshal(t.task)
		executeAt := time.Now().Add(t.delay)
		rdb.ZAdd(ctx, queueKey, redis.Z{
			Score:  float64(executeAt.UnixMilli()),
			Member: string(data),
		})
		fmt.Printf("添加任务: %s, 延迟: %v\n", t.task.ID, t.delay)
	}

	// 消费：轮询取到期任务
	fmt.Println("\n开始消费延迟任务...")
	consumed := 0
	for consumed < 3 {
		now := float64(time.Now().UnixMilli())

		// ZRANGEBYSCORE：获取到期的任务（Score <= 当前时间）
		results, _ := rdb.ZRangeByScore(ctx, queueKey, &redis.ZRangeBy{
			Min:   "-inf",
			Max:   fmt.Sprintf("%f", now),
			Count: 10,
		}).Result()

		for _, raw := range results {
			// ZREM：原子删除（防止多消费者重复消费）
			removed, _ := rdb.ZRem(ctx, queueKey, raw).Result()
			if removed == 0 {
				continue // 已被其他消费者消费
			}

			var task DelayedTask
			json.Unmarshal([]byte(raw), &task)
			fmt.Printf("  执行任务: id=%s, type=%s, payload=%s\n", task.ID, task.Type, task.Payload)
			consumed++
		}

		if consumed < 3 {
			time.Sleep(100 * time.Millisecond) // 轮询间隔
		}
	}

	rdb.Del(ctx, queueKey)
	fmt.Println()
}

// ============================================================
// 2. 生产者-消费者模式（Lua 原子取任务）
// ============================================================

// Lua 脚本：原子地获取并删除到期任务
var fetchTaskScript = redis.NewScript(`
local key = KEYS[1]
local now = ARGV[1]
local limit = ARGV[2]

-- 获取到期任务
local tasks = redis.call('ZRANGEBYSCORE', key, '-inf', now, 'LIMIT', 0, limit)
if #tasks == 0 then
    return {}
end

-- 原子删除已获取的任务
for _, task in ipairs(tasks) do
    redis.call('ZREM', key, task)
end

return tasks
`)

func producerConsumer(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 生产者-消费者模式 ==========")

	queueKey := "delayed:queue:lua"

	// 生产者：持续添加任务
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			task := DelayedTask{
				ID:      fmt.Sprintf("task-%d", i),
				Type:    "order_timeout",
				Payload: fmt.Sprintf("订单 %d 超时取消", 10000+i),
			}
			data, _ := json.Marshal(task)
			// 延迟 100ms ~ 1s
			delay := time.Duration(i*100) * time.Millisecond
			executeAt := time.Now().Add(delay)
			rdb.ZAdd(ctx, queueKey, redis.Z{
				Score:  float64(executeAt.UnixMilli()),
				Member: string(data),
			})
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// 消费者：Lua 脚本原子取任务
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumed := 0
		for {
			select {
			case <-stopCh:
				return
			default:
			}

			now := fmt.Sprintf("%d", time.Now().UnixMilli())
			results, err := fetchTaskScript.Run(ctx, rdb, []string{queueKey}, now, "5").StringSlice()
			if err != nil && err != redis.Nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			for _, raw := range results {
				var task DelayedTask
				json.Unmarshal([]byte(raw), &task)
				fmt.Printf("  [消费者] 执行: %s → %s\n", task.ID, task.Payload)
				consumed++
			}

			if consumed >= 10 {
				return
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()

	wg.Wait()
	close(stopCh)

	// 查看队列剩余
	remaining, _ := rdb.ZCard(ctx, queueKey).Result()
	fmt.Printf("队列剩余任务: %d\n", remaining)

	rdb.Del(ctx, queueKey)
	fmt.Println()

	log.Println("[09_delayed_queue] 延迟队列演示完成")
}
