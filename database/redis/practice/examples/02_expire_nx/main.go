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

// 演示过期时间、NX/XX 条件设置、TTL 查询等企业级常用操作
func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	expireOps(ctx, rdb)
	setConditions(ctx, rdb)
	keyManagement(ctx, rdb)
}

// ============================================================
// 1. 过期时间设置与查询
// ============================================================
func expireOps(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== 过期时间 ==========")

	// SET 时指定过期时间（最推荐的方式，原子操作）
	rdb.Set(ctx, "session:token-abc", `{"user_id":1,"role":"admin"}`, 30*time.Second)

	// TTL 查询剩余过期时间
	ttl, _ := rdb.TTL(ctx, "session:token-abc").Result()
	fmt.Println("TTL session =", ttl) // ~30s

	// PTTL 毫秒精度
	pttl, _ := rdb.PTTL(ctx, "session:token-abc").Result()
	fmt.Println("PTTL session =", pttl)

	// EXPIRE 给已有 key 设置过期时间
	rdb.Set(ctx, "cache:product:1", "iPhone", 0) // 先创建，不过期
	rdb.Expire(ctx, "cache:product:1", 60*time.Second)
	ttl, _ = rdb.TTL(ctx, "cache:product:1").Result()
	fmt.Println("TTL cache:product:1 =", ttl)

	// EXPIREAT 设置绝对过期时间点
	expireAt := time.Now().Add(2 * time.Hour)
	rdb.ExpireAt(ctx, "cache:product:1", expireAt)
	ttl, _ = rdb.TTL(ctx, "cache:product:1").Result()
	fmt.Println("TTL after ExpireAt =", ttl) // ~7200s

	// PERSIST 移除过期时间，变成永久
	rdb.Persist(ctx, "cache:product:1")
	ttl, _ = rdb.TTL(ctx, "cache:product:1").Result()
	fmt.Println("TTL after Persist =", ttl) // -1（永久）

	// 过期时间加随机值（防缓存雪崩）
	baseExpire := 3600 * time.Second
	jitter := time.Duration(time.Now().UnixNano()%300) * time.Second // 0~300秒随机
	rdb.Set(ctx, "cache:hot:1", "data", baseExpire+jitter)
	ttl, _ = rdb.TTL(ctx, "cache:hot:1").Result()
	fmt.Printf("带随机抖动的 TTL = %v (base=3600s + jitter=0~300s)\n", ttl)

	rdb.Del(ctx, "session:token-abc", "cache:product:1", "cache:hot:1")
	fmt.Println()
}

// ============================================================
// 2. SET 的条件参数 NX / XX / GET
// ============================================================
func setConditions(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== SET 条件参数 ==========")

	// NX：不存在才设置（分布式锁的基础）
	ok, _ := rdb.SetNX(ctx, "lock:resource:1", "holder-A", 30*time.Second).Result()
	fmt.Println("SETNX first =", ok) // true

	ok, _ = rdb.SetNX(ctx, "lock:resource:1", "holder-B", 30*time.Second).Result()
	fmt.Println("SETNX second =", ok) // false（已存在，拿不到锁）

	// XX：存在才更新（只更新已有的 key，不创建新 key）
	rdb.SetXX(ctx, "lock:resource:1", "holder-A-renewed", 60*time.Second)
	val, _ := rdb.Get(ctx, "lock:resource:1").Result()
	fmt.Println("After SetXX =", val) // holder-A-renewed

	rdb.SetXX(ctx, "nonexistent", "value", 60*time.Second)
	exists, _ := rdb.Exists(ctx, "nonexistent").Result()
	fmt.Println("SetXX on nonexistent, exists =", exists) // 0（不会创建）

	// SET key value EX seconds NX（生产标准用法：原子加锁）
	// go-redis 用 SetArgs 完整控制所有参数
	result, err := rdb.SetArgs(ctx, "lock:order:99", "uuid-123", redis.SetArgs{
		Mode: "NX",
		TTL:  30 * time.Second,
	}).Result()
	fmt.Println("SetArgs NX result =", result, "err =", err)

	rdb.Del(ctx, "lock:resource:1", "lock:order:99")
	fmt.Println()
}

// ============================================================
// 3. Key 管理命令
// ============================================================
func keyManagement(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Key 管理 ==========")

	rdb.Set(ctx, "test:a", "1", 0)
	rdb.Set(ctx, "test:b", "2", 0)
	rdb.Set(ctx, "test:c", "3", 0)

	// EXISTS 检查 key 是否存在
	count, _ := rdb.Exists(ctx, "test:a", "test:b", "test:notexist").Result()
	fmt.Println("EXISTS count =", count) // 2（只有 a 和 b 存在）

	// TYPE 查看 key 的数据类型
	typ, _ := rdb.Type(ctx, "test:a").Result()
	fmt.Println("TYPE test:a =", typ) // string

	// RENAME 重命名
	rdb.Rename(ctx, "test:c", "test:c_new")
	val, _ := rdb.Get(ctx, "test:c_new").Result()
	fmt.Println("After RENAME, test:c_new =", val)

	// DEL 删除（同步，会阻塞）
	deleted, _ := rdb.Del(ctx, "test:a").Result()
	fmt.Println("DEL test:a, deleted =", deleted)

	// UNLINK 异步删除（大 key 推荐用这个，不阻塞主线程）
	rdb.Unlink(ctx, "test:b", "test:c_new")

	// SCAN 遍历 key（生产环境代替 KEYS 命令，不会阻塞）
	rdb.Set(ctx, "scan:1", "a", 0)
	rdb.Set(ctx, "scan:2", "b", 0)
	rdb.Set(ctx, "scan:3", "c", 0)

	var cursor uint64
	var keys []string
	for {
		var batch []string
		batch, cursor, _ = rdb.Scan(ctx, cursor, "scan:*", 10).Result()
		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}
	fmt.Println("SCAN scan:* =", keys)

	rdb.Del(ctx, "scan:1", "scan:2", "scan:3")
	fmt.Println()

	log.Println("[02_expire_nx] 过期时间与条件设置演示完成")
}
