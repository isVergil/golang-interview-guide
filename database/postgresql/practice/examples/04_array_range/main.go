package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 04_array_range：PG 原生数组 + 范围类型
//
// 数组：适合标签、权限列表等场景，配合 GIN 索引高效查询
// 范围类型：适合时间段、价格区间等场景，配合排除约束防重叠
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	arrayOps(ctx, pool)
	rangeOps(ctx, pool)
}

// ============================================================
// 1. 数组类型
// ============================================================
func arrayOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 数组类型 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS articles`)
	pool.Exec(ctx, `
		CREATE TABLE articles (
			id    serial PRIMARY KEY,
			title text NOT NULL,
			tags  text[] NOT NULL DEFAULT '{}'
		)
	`)

	// 插入数组数据
	pool.Exec(ctx, `INSERT INTO articles (title, tags) VALUES
		('Go 入门教程', ARRAY['go','tutorial','backend']),
		('Redis 缓存设计', ARRAY['redis','cache','backend']),
		('React 实战', ARRAY['react','frontend','javascript']),
		('Go + Redis 分布式锁', ARRAY['go','redis','distributed']),
		('PostgreSQL JSONB', ARRAY['postgresql','jsonb','database'])
	`)
	fmt.Println("插入 5 篇文章")

	// ANY：包含某个元素
	fmt.Println("\n--- 包含 'go' 标签 ---")
	rows, _ := pool.Query(ctx,
		`SELECT title, tags FROM articles WHERE 'go' = ANY(tags)`)
	defer rows.Close()
	for rows.Next() {
		var title string
		var tags []string
		rows.Scan(&title, &tags)
		fmt.Printf("  %s → %v\n", title, tags)
	}

	// @> 包含所有
	fmt.Println("\n--- 同时包含 'go' 和 'redis' ---")
	rows2, _ := pool.Query(ctx,
		`SELECT title FROM articles WHERE tags @> ARRAY['go','redis']`)
	defer rows2.Close()
	for rows2.Next() {
		var title string
		rows2.Scan(&title)
		fmt.Printf("  %s\n", title)
	}

	// && 有交集
	fmt.Println("\n--- 包含 'frontend' 或 'database' ---")
	rows3, _ := pool.Query(ctx,
		`SELECT title FROM articles WHERE tags && ARRAY['frontend','database']`)
	defer rows3.Close()
	for rows3.Next() {
		var title string
		rows3.Scan(&title)
		fmt.Printf("  %s\n", title)
	}

	// array_length：数组长度
	fmt.Println("\n--- 标签数量 ---")
	rows4, _ := pool.Query(ctx,
		`SELECT title, array_length(tags, 1) AS tag_count FROM articles ORDER BY tag_count DESC`)
	defer rows4.Close()
	for rows4.Next() {
		var title string
		var count int
		rows4.Scan(&title, &count)
		fmt.Printf("  %s: %d 个标签\n", title, count)
	}

	// unnest：展开数组（统计每个标签的文章数）
	fmt.Println("\n--- 标签统计（unnest 展开）---")
	rows5, _ := pool.Query(ctx,
		`SELECT tag, COUNT(*) AS cnt FROM articles, unnest(tags) AS tag GROUP BY tag ORDER BY cnt DESC`)
	defer rows5.Close()
	for rows5.Next() {
		var tag string
		var cnt int
		rows5.Scan(&tag, &cnt)
		fmt.Printf("  %s: %d 篇\n", tag, cnt)
	}

	// array_append / array_remove：修改数组
	fmt.Println("\n--- 修改数组 ---")
	pool.Exec(ctx, `UPDATE articles SET tags = array_append(tags, 'hot') WHERE id = 1`)
	pool.Exec(ctx, `UPDATE articles SET tags = array_remove(tags, 'tutorial') WHERE id = 1`)
	var tags []string
	pool.QueryRow(ctx, `SELECT tags FROM articles WHERE id = 1`).Scan(&tags)
	fmt.Printf("  id=1 修改后: %v\n", tags)

	// GIN 索引
	pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_articles_tags ON articles USING GIN (tags)`)
	fmt.Println("  创建 GIN 索引: idx_articles_tags")
	fmt.Println()
}

// ============================================================
// 2. 范围类型
// ============================================================
func rangeOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 范围类型 ==========")

	// 需要 btree_gist 扩展（排除约束需要）
	pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS btree_gist`)

	pool.Exec(ctx, `DROP TABLE IF EXISTS bookings`)
	pool.Exec(ctx, `
		CREATE TABLE bookings (
			id     serial PRIMARY KEY,
			room   text NOT NULL,
			during tstzrange NOT NULL,
			-- 排除约束：组合起来就是，任意两行，如果 room 相等 AND during 重叠，则拒绝。
			EXCLUDE USING GIST (room WITH =, during WITH &&)
		)
	`)
	fmt.Println("创建 bookings 表（带排除约束）")

	// 正常预订
	_, err := pool.Exec(ctx,
		`INSERT INTO bookings (room, during) VALUES ('A101', '[2026-03-25 09:00, 2026-03-25 11:00)')`)
	if err != nil {
		fmt.Printf("预订1失败: %v\n", err)
	} else {
		fmt.Println("预订1: A101 09:00-11:00 成功")
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO bookings (room, during) VALUES ('A101', '[2026-03-25 13:00, 2026-03-25 15:00)')`)
	if err != nil {
		fmt.Printf("预订2失败: %v\n", err)
	} else {
		fmt.Println("预订2: A101 13:00-15:00 成功")
	}

	// 重叠预订 → 自动报错
	_, err = pool.Exec(ctx,
		`INSERT INTO bookings (room, during) VALUES ('A101', '[2026-03-25 10:00, 2026-03-25 12:00)')`)
	if err != nil {
		fmt.Printf("预订3: A101 10:00-12:00 失败（预期，时间重叠）\n")
	}

	// 不同房间不冲突
	_, err = pool.Exec(ctx,
		`INSERT INTO bookings (room, during) VALUES ('B202', '[2026-03-25 09:00, 2026-03-25 11:00)')`)
	if err != nil {
		fmt.Printf("预订4失败: %v\n", err)
	} else {
		fmt.Println("预订4: B202 09:00-11:00 成功（不同房间不冲突）")
	}

	// @> 包含查询：某时间点哪些房间被占用
	fmt.Println("\n--- 10:00 时被占用的房间 ---")
	rows, _ := pool.Query(ctx,
		`SELECT room, during FROM bookings WHERE during @> '2026-03-25 10:00'::timestamptz`)
	defer rows.Close()
	for rows.Next() {
		var room, during string
		rows.Scan(&room, &during)
		fmt.Printf("  %s: %s\n", room, during)
	}

	// && 重叠查询：某时间段内有哪些预订
	fmt.Println("\n--- 08:00-12:00 期间的所有预订 ---")
	rows2, _ := pool.Query(ctx,
		`SELECT room, during FROM bookings 
		 WHERE during && '[2026-03-25 08:00, 2026-03-25 12:00)'::tstzrange`)
	defer rows2.Close()
	for rows2.Next() {
		var room, during string
		rows2.Scan(&room, &during)
		fmt.Printf("  %s: %s\n", room, during)
	}

	log.Println("[04_array_range] 数组和范围类型演示完成")
}
