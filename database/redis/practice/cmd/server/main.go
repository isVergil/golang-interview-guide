package main

import (
	"context"
	"fmt"
	"log"

	"redis-practice/internal/config"
	pkgRedis "redis-practice/pkg/redis"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	fmt.Println("[1] 配置加载成功")

	// 2. 初始化 Redis 连接
	rdb, err := pkgRedis.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("init redis: %v", err)
	}
	defer pkgRedis.Close(rdb)
	fmt.Println("[2] Redis 连接成功")

	// 3. 验证连接信息
	ctx := context.Background()
	info, _ := rdb.Info(ctx, "server").Result()
	_ = info // 连接正常即可
	dbSize, _ := rdb.DBSize(ctx).Result()
	fmt.Printf("[3] 当前 DB key 数量: %d\n", dbSize)

	fmt.Println("\n基建就绪，可运行 examples/ 下的各个示例")
	fmt.Println("示例列表:")
	fmt.Println("  01_basic_types       - 五种基本数据类型")
	fmt.Println("  02_expire_nx         - 过期策略与条件操作")
	fmt.Println("  03_pipeline          - Pipeline 批量操作")
	fmt.Println("  04_lua               - Lua 脚本原子操作")
	fmt.Println("  05_distributed_lock  - 分布式锁")
	fmt.Println("  06_cache_patterns    - 缓存模式（穿透/击穿/旁路）")
	fmt.Println("  07_bitmap_hyperloglog - Bitmap 签到 & HyperLogLog UV")
	fmt.Println("  08_pubsub            - 发布/订阅")
	fmt.Println("  09_delayed_queue     - ZSet 延迟队列")
	fmt.Println("  10_rate_limiter      - 限流器（固定窗口/滑动窗口/令牌桶）")
}
