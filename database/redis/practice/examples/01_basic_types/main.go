package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

// 演示 Redis 五大基础数据类型的 CRUD 操作
func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	stringOps(ctx, rdb)
	hashOps(ctx, rdb)
	listOps(ctx, rdb)
	setOps(ctx, rdb)
	zsetOps(ctx, rdb)
}

// ============================================================
// 1. String：最基础的 key-value
// ============================================================
func stringOps(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== String ==========")

	// SET / GET
	rdb.Set(ctx, "user:name", "alice", 0)
	val, _ := rdb.Get(ctx, "user:name").Result()
	fmt.Println("GET user:name =", val)

	// INCR 原子自增（计数器场景）
	rdb.Set(ctx, "page:views", 0, 0)
	rdb.Incr(ctx, "page:views")
	rdb.IncrBy(ctx, "page:views", 10)
	views, _ := rdb.Get(ctx, "page:views").Int64()
	fmt.Println("page:views =", views) // 11

	// SETNX 不存在才设置（分布式锁基础）
	ok, _ := rdb.SetNX(ctx, "lock:order:1", "uuid-xxx", 0).Result()
	fmt.Println("SETNX lock:order:1 =", ok) // true（首次）
	ok, _ = rdb.SetNX(ctx, "lock:order:1", "uuid-yyy", 0).Result()
	fmt.Println("SETNX lock:order:1 =", ok) // false（已存在）

	// MSET / MGET 批量操作（减少网络往返）
	rdb.MSet(ctx, "k1", "v1", "k2", "v2", "k3", "v3")
	vals, _ := rdb.MGet(ctx, "k1", "k2", "k3").Result()
	fmt.Println("MGET =", vals)

	// 查看底层编码
	rdb.Set(ctx, "num", 123, 0)
	enc, _ := rdb.ObjectEncoding(ctx, "num").Result()
	fmt.Println("num encoding =", enc) // int

	rdb.Set(ctx, "short", "hello", 0)
	enc, _ = rdb.ObjectEncoding(ctx, "short").Result()
	fmt.Println("short encoding =", enc) // embstr

	// 清理
	rdb.Del(ctx, "user:name", "page:views", "lock:order:1", "k1", "k2", "k3", "num", "short")
	fmt.Println()
}

// ============================================================
// 2. Hash：field-value 映射，适合存对象
// ============================================================
func hashOps(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Hash ==========")

	// HSET 批量设置字段
	rdb.HSet(ctx, "user:1", map[string]interface{}{
		"name":  "alice",
		"age":   25,
		"email": "alice@test.com",
	})

	// HGET 取单个字段
	name, _ := rdb.HGet(ctx, "user:1", "name").Result()
	fmt.Println("HGET name =", name)

	// HGETALL 取全部字段
	all, _ := rdb.HGetAll(ctx, "user:1").Result()
	fmt.Println("HGETALL =", all)

	// HINCRBY 单字段自增（购物车加数量）
	rdb.HIncrBy(ctx, "user:1", "age", 1)
	age, _ := rdb.HGet(ctx, "user:1", "age").Result()
	fmt.Println("HINCRBY age =", age) // 26

	// HDEL 删除单个字段
	//rdb.HDel(ctx, "user:1", "email")

	// HLEN 字段数量
	length, _ := rdb.HLen(ctx, "user:1").Result()
	fmt.Println("HLEN =", length) // 2

	// HEXISTS 判断字段是否存在
	exists, _ := rdb.HExists(ctx, "user:1", "email").Result()
	fmt.Println("HEXISTS email =", exists) // false

	// 查看底层编码
	enc, _ := rdb.ObjectEncoding(ctx, "user:1").Result()
	fmt.Println("user:1 encoding =", enc) // listpack（字段少时）

	rdb.Del(ctx, "user:1")
	fmt.Println()
}

// ============================================================
// 3. List：有序列表，支持头尾操作
// ============================================================
func listOps(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== List ==========")

	// LPUSH 左侧插入（头部），RPUSH 右侧插入（尾部）
	rdb.RPush(ctx, "queue", "task1", "task2", "task3")
	rdb.LPush(ctx, "queue", "task0") // 插到最前面

	// LRANGE 获取范围（0 到 -1 = 全部）
	items, _ := rdb.LRange(ctx, "queue", 0, -1).Result()
	fmt.Println("LRANGE queue =", items) // [task0, task1, task2, task3]

	// LPOP / RPOP 弹出
	left, _ := rdb.LPop(ctx, "queue").Result()
	fmt.Println("LPOP =", left) // task0
	right, _ := rdb.RPop(ctx, "queue").Result()
	fmt.Println("RPOP =", right) // task3

	// LLEN 长度
	length, _ := rdb.LLen(ctx, "queue").Result()
	fmt.Println("LLEN =", length) // 2

	// LINDEX 按索引取值
	val, _ := rdb.LIndex(ctx, "queue", 0).Result()
	fmt.Println("LINDEX 0 =", val) // task1

	rdb.Del(ctx, "queue")
	fmt.Println()
}

// ============================================================
// 4. Set：无序不重复集合
// ============================================================
func setOps(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Set ==========")

	// SADD 添加成员
	rdb.SAdd(ctx, "tags:article:1", "go", "redis", "mysql")
	rdb.SAdd(ctx, "tags:article:2", "go", "kafka", "docker")

	// SMEMBERS 获取所有成员
	members, _ := rdb.SMembers(ctx, "tags:article:1").Result()
	fmt.Println("SMEMBERS article:1 =", members)

	// SISMEMBER 判断是否存在
	isMember, _ := rdb.SIsMember(ctx, "tags:article:1", "go").Result()
	fmt.Println("SISMEMBER go =", isMember) // true

	// SCARD 成员数量
	count, _ := rdb.SCard(ctx, "tags:article:1").Result()
	fmt.Println("SCARD =", count) // 3

	// SINTER 交集（共同标签）
	inter, _ := rdb.SInter(ctx, "tags:article:1", "tags:article:2").Result()
	fmt.Println("SINTER =", inter) // [go]

	// SUNION 并集
	union, _ := rdb.SUnion(ctx, "tags:article:1", "tags:article:2").Result()
	fmt.Println("SUNION =", union)

	// SDIFF 差集（在 A 中但不在 B 中，顺序有关）
	// article:1 = {go, redis, mysql}
	// article:2 = {go, kafka, docker}
	// SDIFF article:1 article:2 → {redis, mysql}（go 是交集，排除掉）
	diff, _ := rdb.SDiff(ctx, "tags:article:1", "tags:article:2").Result()
	fmt.Println("SDIFF =", diff) // [redis, mysql]

	// SRANDMEMBER 随机取（抽奖场景）
	random, _ := rdb.SRandMemberN(ctx, "tags:article:1", 2).Result()
	fmt.Println("SRANDMEMBER 2 =", random)

	rdb.Del(ctx, "tags:article:1", "tags:article:2")
	fmt.Println()
}

// ============================================================
// 5. ZSet（Sorted Set）：有序集合，按 score 排序
// ============================================================
func zsetOps(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== ZSet ==========")

	// ZADD 添加成员和分数
	rdb.ZAdd(ctx, "ranking", redis.Z{Score: 100, Member: "alice"})
	rdb.ZAdd(ctx, "ranking", redis.Z{Score: 90, Member: "bob"})
	rdb.ZAdd(ctx, "ranking", redis.Z{Score: 95, Member: "charlie"})
	rdb.ZAdd(ctx, "ranking", redis.Z{Score: 80, Member: "dave"})

	// ZRANGE 按 score 正序（参数是索引范围，0 到 -1 表示全部）
	members, _ := rdb.ZRangeWithScores(ctx, "ranking", 0, -1).Result()
	fmt.Println("ZRANGE (正序):")
	for _, m := range members {
		fmt.Printf("  %s: %.0f\n", m.Member, m.Score)
	}

	// ZREVRANGE 按 score 倒序（0 到 2 = 前 3 名，索引从 0 开始）
	top, _ := rdb.ZRevRangeWithScores(ctx, "ranking", 0, 2).Result()
	fmt.Println("Top 3 (倒序):")
	for i, m := range top {
		fmt.Printf("  #%d %s: %.0f\n", i+1, m.Member, m.Score)
	}

	// ZSCORE 查分数
	score, _ := rdb.ZScore(ctx, "ranking", "alice").Result()
	fmt.Println("ZSCORE alice =", score) // 100

	// ZRANK 查排名（正序，从 0 开始）
	rank, _ := rdb.ZRank(ctx, "ranking", "alice").Result()
	fmt.Println("ZRANK alice =", rank) // 3（正序第4）

	// ZREVRANK 查倒序排名
	revRank, _ := rdb.ZRevRank(ctx, "ranking", "alice").Result()
	fmt.Println("ZREVRANK alice =", revRank) // 0（倒序第1）

	// ZINCRBY 加分
	newScore, _ := rdb.ZIncrBy(ctx, "ranking", 15, "bob").Result()
	fmt.Println("ZINCRBY bob +15 =", newScore) // 105

	// ZRANGEBYSCORE 按分数范围查（默认闭区间 >=Min && <=Max）
	// 开区间用 "(" 前缀：Min: "(90" 表示 >90，Max: "(100" 表示 <100
	byScore, _ := rdb.ZRangeByScoreWithScores(ctx, "ranking", &redis.ZRangeBy{
		Min: "90",  // >= 90
		Max: "100", // <= 100
	}).Result()
	fmt.Println("ZRANGEBYSCORE 90~100:")
	for _, m := range byScore {
		fmt.Printf("  %s: %.0f\n", m.Member, m.Score)
	}

	// ZCARD 成员总数
	total, _ := rdb.ZCard(ctx, "ranking").Result()
	fmt.Println("ZCARD =", total)

	// 查看底层编码
	enc, _ := rdb.ObjectEncoding(ctx, "ranking").Result()
	fmt.Printf("ranking encoding = %s (元素少时用 listpack，多时用 skiplist)\n", enc)

	rdb.Del(ctx, "ranking")
	fmt.Println()

	log.Println("[01_basic_types] 五大数据类型演示完成")
}
