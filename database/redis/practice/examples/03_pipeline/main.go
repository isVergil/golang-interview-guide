package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

// 演示 Pipeline 批量操作，减少网络往返
func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	pipelineBasic(ctx, rdb)
	pipelineTransaction(ctx, rdb)
	pipelineBenchmark(ctx, rdb)
}

// ============================================================
// 1. 基础 Pipeline：多个命令打包发送
// ============================================================
func pipelineBasic(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Pipeline 基础 ==========")

	// Pipeline 把多条命令打包一次发送，减少 RTT
	pipe := rdb.Pipeline()

	// 在 Pipeline 里塞命令（不会立即执行）
	incr := pipe.Incr(ctx, "pipeline:counter")
	pipe.Set(ctx, "pipeline:name", "alice", time.Minute)
	pipe.HSet(ctx, "pipeline:user", "age", 25)
	get := pipe.Get(ctx, "pipeline:name")

	// Exec 一次性发送所有命令
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("pipeline exec: %v", err)
		return
	}

	// 读取每个命令的结果
	fmt.Println("INCR counter =", incr.Val()) // 1
	fmt.Println("GET name =", get.Val())      // alice

	rdb.Del(ctx, "pipeline:counter", "pipeline:name", "pipeline:user")
	fmt.Println()
}

// ============================================================
// 2. TxPipeline：带 MULTI/EXEC 事务的 Pipeline
// ============================================================
func pipelineTransaction(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== TxPipeline（事务） ==========")

	// TxPipeline = Pipeline + MULTI/EXEC 包裹
	// 保证这批命令要么全部执行，要么全不执行（Redis 层面原子）
	pipe := rdb.TxPipeline()

	pipe.Set(ctx, "tx:balance:A", 100, 0)
	pipe.Set(ctx, "tx:balance:B", 200, 0)
	pipe.DecrBy(ctx, "tx:balance:A", 30) // A 扣 30
	pipe.IncrBy(ctx, "tx:balance:B", 30) // B 加 30

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("tx pipeline exec: %v", err)
		return
	}

	a, _ := rdb.Get(ctx, "tx:balance:A").Int64()
	b, _ := rdb.Get(ctx, "tx:balance:B").Int64()
	fmt.Printf("转账后：A=%d, B=%d (总和=%d，不变)\n", a, b, a+b) // 70, 230, 300

	rdb.Del(ctx, "tx:balance:A", "tx:balance:B")
	fmt.Println()
}

// ============================================================
// 3. Pipeline vs 逐条执行 性能对比
// ============================================================
func pipelineBenchmark(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Pipeline 性能对比 ==========")
	n := 1000

	// 逐条执行
	start := time.Now()
	for i := 0; i < n; i++ {
		rdb.Set(ctx, fmt.Sprintf("bench:single:%d", i), i, time.Minute)
	}
	singleDur := time.Since(start)
	fmt.Printf("逐条执行 %d 次 SET: %v\n", n, singleDur)

	// Pipeline 批量执行
	start = time.Now()
	pipe := rdb.Pipeline()
	for i := 0; i < n; i++ {
		pipe.Set(ctx, fmt.Sprintf("bench:pipe:%d", i), i, time.Minute)
	}
	pipe.Exec(ctx)
	pipeDur := time.Since(start)
	fmt.Printf("Pipeline  %d 次 SET: %v\n", n, pipeDur)

	if singleDur > 0 {
		fmt.Printf("Pipeline 快了约 %.1f 倍\n", float64(singleDur)/float64(pipeDur))
	}

	// 清理
	pipe = rdb.Pipeline()
	for i := 0; i < n; i++ {
		pipe.Del(ctx, fmt.Sprintf("bench:single:%d", i))
		pipe.Del(ctx, fmt.Sprintf("bench:pipe:%d", i))
	}
	pipe.Exec(ctx)
	fmt.Println()

	log.Println("[03_pipeline] Pipeline 演示完成")
}
