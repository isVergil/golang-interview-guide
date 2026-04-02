package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	basicPubSub(ctx, rdb)
	patternSubscribe(ctx, rdb)
	multiSubscribers(ctx, rdb)
}

// ============================================================
// 1. 基础 Pub/Sub
//
// 即发即弃（Fire and Forget）：
//   - 消息不持久化，Redis 不保存历史消息
//   - 订阅者离线期间的消息永久丢失，重连后收不到
//   - 没有 ACK 机制，Redis 不知道订阅者是否成功处理
//   - 需要可靠消息用 Redis Stream（5.0+）或 Kafka
//
// PUBLISH 性能与订阅者数量：
//   - PUBLISH 是同步操作，Redis 主线程逐个给每个订阅者发送消息
//   - 订阅者越多，PUBLISH 耗时越长（O(N)，N=订阅者数量）
//   - 如果某个订阅者的输出缓冲区满了（消费太慢），Redis 会断开该连接
//   - 生产环境建议控制单频道订阅者数量，大规模广播用 Kafka
//
// 使用场景：
//   - 实时通知：聊天室消息、系统公告、配置变更通知
//   - 事件驱动：订单状态变更通知下游服务（对丢失不敏感）
//   - 缓存失效：一个服务更新数据后通知其他服务清除本地缓存
//   - 实时监控：指标数据实时推送到 Dashboard
//
// ============================================================
func basicPubSub(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 基础 Pub/Sub ==========")

	channel := "chat:room:1"

	// 订阅（必须在发布之前）
	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()

	// 等待订阅确认
	_, err := sub.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("已订阅频道:", channel)

	// 消费者：异步接收消息
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := sub.Channel()
		count := 0
		for msg := range ch {
			fmt.Printf("  收到消息: channel=%s, payload=%s\n", msg.Channel, msg.Payload)
			count++
			if count >= 3 {
				return
			}
		}
	}()

	// 生产者：发布消息
	time.Sleep(50 * time.Millisecond) // 确保订阅者就绪
	for i := 1; i <= 3; i++ {
		receivers, _ := rdb.Publish(ctx, channel, fmt.Sprintf("hello #%d", i)).Result()
		fmt.Printf("  发布消息 #%d, 接收者数量: %d\n", i, receivers)
		time.Sleep(20 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println()
}

// ============================================================
// 2. 模式订阅（Pattern Subscribe）
// ============================================================
func patternSubscribe(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 模式订阅 PSubscribe ==========")

	// 订阅所有 order:* 频道
	sub := rdb.PSubscribe(ctx, "order:*")
	defer sub.Close()

	_, err := sub.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("已订阅模式: order:*")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := sub.Channel()
		count := 0
		for msg := range ch {
			fmt.Printf("  模式匹配: pattern=%s, channel=%s, payload=%s\n",
				msg.Pattern, msg.Channel, msg.Payload)
			count++
			if count >= 3 {
				return
			}
		}
	}()

	time.Sleep(50 * time.Millisecond)
	// 发布到不同的 order 子频道
	rdb.Publish(ctx, "order:created", `{"id":1001,"status":"created"}`)
	time.Sleep(10 * time.Millisecond)
	rdb.Publish(ctx, "order:paid", `{"id":1001,"status":"paid"}`)
	time.Sleep(10 * time.Millisecond)
	rdb.Publish(ctx, "order:shipped", `{"id":1001,"status":"shipped"}`)

	wg.Wait()
	fmt.Println()
}

// ============================================================
// 3. 多订阅者（广播模式）
// ============================================================
func multiSubscribers(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 多订阅者（广播）==========")

	channel := "broadcast:notify"

	var wg sync.WaitGroup

	// 启动 3 个订阅者
	for i := 1; i <= 3; i++ {
		sub := rdb.Subscribe(ctx, channel)
		_, err := sub.Receive(ctx)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func(id int, s *redis.PubSub) {
			defer wg.Done()
			defer s.Close()
			ch := s.Channel()
			for msg := range ch {
				fmt.Printf("  订阅者 #%d 收到: %s\n", id, msg.Payload)
				return // 每个订阅者只读一条
			}
		}(i, sub)
	}

	time.Sleep(50 * time.Millisecond)

	// 发布一条消息 → 所有订阅者都会收到（广播）
	receivers, _ := rdb.Publish(ctx, channel, "系统维护通知：今晚 22:00 维护").Result()
	fmt.Printf("  发布广播，接收者: %d\n", receivers)

	wg.Wait()
	fmt.Println()

	log.Println("[08_pubsub] Pub/Sub 演示完成")
}
