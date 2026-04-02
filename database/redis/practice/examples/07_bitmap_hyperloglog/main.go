package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"redis-practice/examples/setup"
	pkgRedis "redis-practice/pkg/redis"
)

func main() {
	rdb := setup.MustSetup()
	defer pkgRedis.Close(rdb)
	ctx := context.Background()

	bitmapSignIn(ctx, rdb)
	bitmapStats(ctx, rdb)
	hyperLogLogUV(ctx, rdb)
}

// ============================================================
// 1. Bitmap 签到系统
// ============================================================
func bitmapSignIn(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Bitmap 签到系统 ==========")

	userID := 1001
	now := time.Now()
	key := fmt.Sprintf("sign:%d:%s", userID, now.Format("200601"))
	day := now.Day()

	// SETBIT：签到（offset = 日期-1）
	rdb.SetBit(ctx, key, int64(day-1), 1)
	fmt.Printf("用户 %d 在第 %d 天签到成功\n", userID, day)

	// 补签昨天
	rdb.SetBit(ctx, key, int64(day-2), 1)
	fmt.Printf("补签第 %d 天\n", day-1)

	// GETBIT：查询某天是否签到
	signed, _ := rdb.GetBit(ctx, key, int64(day-1)).Result()
	fmt.Printf("第 %d 天是否签到: %d (1=是, 0=否)\n", day, signed)

	// BITCOUNT：统计本月签到天数
	count, _ := rdb.BitCount(ctx, key, &redis.BitCount{Start: 0, End: -1}).Result()
	fmt.Printf("本月累计签到: %d 天\n", count)

	// BITPOS：查找本月第一次签到是第几天
	firstDay, _ := rdb.BitPos(ctx, key, 1).Result()
	fmt.Printf("本月首次签到: 第 %d 天\n", firstDay+1)

	rdb.Del(ctx, key)
	fmt.Println()
}

// ============================================================
// 2. Bitmap 统计：连续签到
// ============================================================
func bitmapStats(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== Bitmap 连续签到统计 ==========")

	key := "sign:demo:202603"
	// 模拟签到数据：第1-5天连续签到, 第7天签到, 第10-15天连续签到
	signDays := []int{1, 2, 3, 4, 5, 7, 10, 11, 12, 13, 14, 15}
	for _, d := range signDays {
		rdb.SetBit(ctx, key, int64(d-1), 1)
	}

	total, _ := rdb.BitCount(ctx, key, &redis.BitCount{Start: 0, End: -1}).Result()
	fmt.Printf("总签到天数: %d\n", total)

	// 获取原始数据来计算连续签到
	bytes, _ := rdb.Get(ctx, key).Bytes()
	maxStreak, currentStreak := 0, 0
	for _, b := range bytes {
		for i := 7; i >= 0; i-- {
			if (b>>uint(i))&1 == 1 {
				currentStreak++
				if currentStreak > maxStreak {
					maxStreak = currentStreak
				}
			} else {
				currentStreak = 0
			}
		}
	}
	fmt.Printf("最长连续签到: %d 天\n", maxStreak)

	// BITOP：多 key 运算
	// 统计两个用户都签到的天数（交集）
	key2 := "sign:demo2:202603"
	for _, d := range []int{3, 4, 5, 6, 7} {
		rdb.SetBit(ctx, key2, int64(d-1), 1)
	}

	destKey := "sign:and:result"
	rdb.BitOpAnd(ctx, destKey, key, key2)
	bothCount, _ := rdb.BitCount(ctx, destKey, &redis.BitCount{Start: 0, End: -1}).Result()
	fmt.Printf("两用户同时签到天数: %d\n", bothCount)

	rdb.Del(ctx, key, key2, destKey)
	fmt.Println()
}

// ============================================================
// 3. HyperLogLog UV 统计
// ============================================================
func hyperLogLogUV(ctx context.Context, rdb *redis.Client) {
	fmt.Println("========== HyperLogLog UV 统计 ==========")

	// 模拟页面 PV/UV 统计
	pageKey := "uv:page:home:20260324"

	// 添加访客（模拟 1000 个 UV，其中有重复）
	pipe := rdb.Pipeline()
	for i := 0; i < 1500; i++ {
		uid := i % 1000 // 1000 个独立用户，500 个重复
		pipe.PFAdd(ctx, pageKey, "user:"+strconv.Itoa(uid))
	}
	pipe.Exec(ctx)

	// PFCOUNT：统计 UV
	uv, _ := rdb.PFCount(ctx, pageKey).Result()
	fmt.Printf("页面 UV（实际 1000，HyperLogLog 估算）: %d\n", uv)
	fmt.Printf("误差率: %.2f%%\n", float64(uv-1000)/1000*100)

	// 模拟多个页面，统计全站 UV
	page2Key := "uv:page:detail:20260324"
	pipe = rdb.Pipeline()
	for i := 500; i < 2000; i++ { // 和 page1 有重叠用户
		pipe.PFAdd(ctx, page2Key, "user:"+strconv.Itoa(i))
	}
	pipe.Exec(ctx)

	uv2, _ := rdb.PFCount(ctx, page2Key).Result()
	fmt.Printf("详情页 UV: %d\n", uv2)

	// PFMERGE：合并统计全站 UV（自动去重）
	mergedKey := "uv:site:20260324"
	rdb.PFMerge(ctx, mergedKey, pageKey, page2Key)
	totalUV, _ := rdb.PFCount(ctx, mergedKey).Result()
	fmt.Printf("全站 UV（合并去重）: %d（实际应约 2000）\n", totalUV)

	// 内存对比
	bitmapKey := "uv:bitmap:demo"
	for i := 0; i < 1000000; i++ {
		rdb.SetBit(ctx, bitmapKey, int64(i), 1)
	}
	bitmapMem, _ := rdb.MemoryUsage(ctx, bitmapKey).Result()
	rdb.Del(ctx, bitmapKey)

	hllKey := "uv:hll:demo"
	pipe = rdb.Pipeline()
	for i := 0; i < 1000000; i++ {
		pipe.PFAdd(ctx, hllKey, "u:"+strconv.Itoa(i))
	}
	pipe.Exec(ctx)
	hllMem, _ := rdb.MemoryUsage(ctx, hllKey).Result()

	fmt.Printf("\n--- 内存对比（100万元素）---\n")
	fmt.Printf("Bitmap 内存: %d bytes (%.2f KB)\n", bitmapMem, float64(bitmapMem)/1024)
	fmt.Printf("HyperLogLog 内存: %d bytes (%.2f KB)\n", hllMem, float64(hllMem)/1024)
	fmt.Printf("HyperLogLog 始终 ≈12 KB（与数据量无关）\n")

	rdb.Del(ctx, pageKey, page2Key, mergedKey, hllKey)
	fmt.Println()

	log.Println("[07_bitmap_hyperloglog] 演示完成")
}
