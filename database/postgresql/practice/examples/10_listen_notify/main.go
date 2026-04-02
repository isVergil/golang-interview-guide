package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 10_listen_notify：PostgreSQL 原生消息通知（Pub/Sub）
//
// LISTEN/NOTIFY 是 PG 内置的轻量级消息通知机制：
//   - LISTEN channel_name   — 订阅频道
//   - NOTIFY channel_name, 'payload' — 发送消息
//   - 消息不持久化（fire-and-forget），未监听就丢失
//   - 多个 listener 可以订阅同一频道
//
// 适合场景：
//   - 缓存失效通知（数据变更 → 通知应用清缓存）
//   - 任务队列（新任务入库 → 通知 worker 处理）
//   - 配置变更广播
//
// 不适合场景：
//   - 需要持久化、重试、确认的消息（用 Kafka/RabbitMQ）
//   - 高吞吐量消息（NOTIFY 走 WAL，性能有上限）
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	initData(ctx, pool)
	basicNotify(ctx, pool)
	triggerNotify(ctx, pool)

	log.Println("[10_listen_notify] 演示完成")
}

func initData(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 初始化 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS orders`)
	pool.Exec(ctx, `
		CREATE TABLE orders (
			id         bigserial PRIMARY KEY,
			user_id    bigint NOT NULL,
			amount     numeric(10,2) NOT NULL,
			status     text NOT NULL DEFAULT 'pending',
			created_at timestamptz NOT NULL DEFAULT now()
		)
	`)

	// 创建触发器：订单表变更时自动 NOTIFY
	pool.Exec(ctx, `
		CREATE OR REPLACE FUNCTION notify_order_change() RETURNS trigger AS $$
		BEGIN
			PERFORM pg_notify('order_changes', json_build_object(
				'op', TG_OP,
				'id', COALESCE(NEW.id, OLD.id),
				'status', COALESCE(NEW.status, OLD.status)
			)::text);
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql
	`)

	pool.Exec(ctx, `DROP TRIGGER IF EXISTS trg_order_notify ON orders`)
	pool.Exec(ctx, `
		CREATE TRIGGER trg_order_notify
		AFTER INSERT OR UPDATE ON orders
		FOR EACH ROW EXECUTE FUNCTION notify_order_change()
	`)
	fmt.Println("创建 orders 表 + 变更通知触发器")
}

// ============================================================
// 1. 基本 LISTEN/NOTIFY
// ============================================================
func basicNotify(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 基本 LISTEN/NOTIFY ==========")

	// LISTEN 需要用独占连接（不能用连接池的共享连接）
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Printf("acquire conn: %v", err)
		return
	}
	defer conn.Release()

	// 订阅频道
	_, err = conn.Exec(ctx, `LISTEN test_channel`)
	if err != nil {
		log.Printf("listen: %v", err)
		return
	}
	fmt.Println("订阅频道: test_channel")

	// 另一个连接发送通知
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)

		// 发送 3 条消息
		for i := 1; i <= 3; i++ {
			payload := fmt.Sprintf(`{"msg":"hello","seq":%d}`, i)
			pool.Exec(ctx, `SELECT pg_notify($1, $2)`, "test_channel", payload)
			fmt.Printf("  发送: %s\n", payload)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 接收通知
	fmt.Println("等待接收消息...")
	received := 0
	timeout := time.After(2 * time.Second)

	for received < 3 {
		select {
		case <-timeout:
			fmt.Println("超时退出")
			wg.Wait()
			return
		default:
			notification, err := conn.Conn().WaitForNotification(ctx)
			if err != nil {
				continue
			}
			received++
			fmt.Printf("  收到 [%s]: %s (PID=%d)\n",
				notification.Channel, notification.Payload, notification.PID)
		}
	}

	// 取消订阅
	conn.Exec(ctx, `UNLISTEN test_channel`)
	fmt.Println("取消订阅 test_channel")
	wg.Wait()
}

// ============================================================
// 2. 触发器 + NOTIFY（实战场景）
// ============================================================
func triggerNotify(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 触发器 + NOTIFY ==========")
	fmt.Println("场景：订单表变更自动通知，实现缓存失效/实时推送")

	// 监听订单变更
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Printf("acquire conn: %v", err)
		return
	}
	defer conn.Release()

	conn.Exec(ctx, `LISTEN order_changes`)
	fmt.Println("订阅频道: order_changes")

	// 启动消费者
	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				// 带超时等待通知
				waitCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
				notification, err := conn.Conn().WaitForNotification(waitCtx)
				cancel()
				if err != nil {
					continue
				}

				// 解析通知负载
				var payload map[string]interface{}
				json.Unmarshal([]byte(notification.Payload), &payload)
				fmt.Printf("  收到变更通知: op=%s, id=%.0f, status=%s\n",
					payload["op"], payload["id"], payload["status"])
			}
		}
	}()

	// 模拟订单操作（会触发 NOTIFY）
	time.Sleep(200 * time.Millisecond)
	fmt.Println("\n模拟订单操作:")

	// 创建订单
	var orderID int64
	pool.QueryRow(ctx,
		`INSERT INTO orders (user_id, amount) VALUES ($1, $2) RETURNING id`,
		1001, 99.99,
	).Scan(&orderID)
	fmt.Printf("  创建订单 id=%d\n", orderID)
	time.Sleep(300 * time.Millisecond)

	// 更新订单状态
	pool.Exec(ctx, `UPDATE orders SET status = 'paid' WHERE id = $1`, orderID)
	fmt.Printf("  订单 %d 状态更新为 paid\n", orderID)
	time.Sleep(300 * time.Millisecond)

	// 再创建一个
	pool.QueryRow(ctx,
		`INSERT INTO orders (user_id, amount) VALUES ($1, $2) RETURNING id`,
		1002, 199.50,
	).Scan(&orderID)
	fmt.Printf("  创建订单 id=%d\n", orderID)
	time.Sleep(300 * time.Millisecond)

	close(done)
	wg.Wait()

	conn.Exec(ctx, `UNLISTEN order_changes`)
	fmt.Println("\n取消订阅，演示结束")

	fmt.Println(`
LISTEN/NOTIFY vs 外部消息队列:
  ┌──────────────┬──────────────┬────────────────┐
  │              │ PG NOTIFY    │ Kafka/RabbitMQ │
  ├──────────────┼──────────────┼────────────────┤
  │ 依赖         │ 无额外组件    │ 需要独立部署    │
  │ 持久化       │ 不持久化      │ 持久化          │
  │ 吞吐量       │ 中等          │ 高             │
  │ 消息确认     │ 无            │ ACK 机制        │
  │ 适合场景     │ 缓存失效通知  │ 解耦/削峰/可靠  │
  └──────────────┴──────────────┴────────────────┘`)
}
