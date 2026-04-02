package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 03_jsonb：JSONB 半结构化存储 + GIN 索引
//
// PostgreSQL 的 JSONB 是二进制 JSON，支持索引和丰富的查询运算符
// 适合日志、事件、配置等 schema 不固定的数据
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	initTable(ctx, pool)
	insertOps(ctx, pool)
	queryOps(ctx, pool)
	updateOps(ctx, pool)
	ginIndex(ctx, pool)
}

func initTable(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 建表 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS events`)
	pool.Exec(ctx, `
		CREATE TABLE events (
			id   bigserial PRIMARY KEY,
			data jsonb NOT NULL,
			created_at timestamptz NOT NULL DEFAULT now()
		)
	`)
	fmt.Println("创建 events 表（data 列为 JSONB）")
}

// ============================================================
// 1. 插入 JSONB 数据
// ============================================================
func insertOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== INSERT JSONB ==========")

	events := []string{
		`{"type":"click","page":"/home","user_id":1001,"meta":{"browser":"chrome","os":"mac"}}`,
		`{"type":"click","page":"/product/1","user_id":1002,"meta":{"browser":"firefox","os":"windows"}}`,
		`{"type":"purchase","page":"/checkout","user_id":1001,"amount":9900,"items":["book","pen"]}`,
		`{"type":"login","user_id":1003,"ip":"192.168.1.100"}`,
		`{"type":"click","page":"/home","user_id":1003,"meta":{"browser":"safari","os":"ios"}}`,
	}

	for _, e := range events {
		var id int64
		var ctid string
		pool.QueryRow(ctx,
			`INSERT INTO events (data) VALUES ($1::jsonb) RETURNING id,ctid`, e,
		).Scan(&id, &ctid)
		fmt.Printf("  插入事件 id=%d, ctid=%s\n", id, ctid)
	}
}

// ============================================================
// 2. JSONB 查询
// ============================================================
func queryOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== JSONB 查询 ==========")

	// ->> 取出的是纯文本字符串，不能再链式取
	fmt.Println("--- 取文本值 (->>) ---")
	rows, _ := pool.Query(ctx,
		`SELECT id, data->>'type' AS event_type, data->>'page' AS page FROM events`)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var eventType, page *string
		rows.Scan(&id, &eventType, &page)
		p := "<nil>"
		if page != nil {
			p = *page
		}
		fmt.Printf("  id=%d, type=%s, page=%s\n", id, *eventType, p)
	}

	// -> 取 JSON 值（嵌套访问），取出的还是 JSON 对象，可以继续链式取值
	fmt.Println("\n--- 嵌套访问 (-> ->>) ---")
	rows2, _ := pool.Query(ctx,
		`SELECT id, data->'meta'->>'browser' AS browser 
		 FROM events WHERE data ? 'meta'`)
	defer rows2.Close()
	for rows2.Next() {
		var id int64
		var browser string
		rows2.Scan(&id, &browser)
		fmt.Printf("  id=%d, browser=%s\n", id, browser)
	}

	// @> 包含查询（最常用）匹配规则是子集匹配，只要 data 包含右边所有的键值对就行，data 里有其他字段不影响匹配
	// ✓ {"type":"click", "page":"/home", "user_id":1}    -- 包含，匹配
	// ✓ {"type":"click", "page":"/home", "extra":"abc"}   -- 包含，匹配
	// ✗ {"type":"click", "page":"/about"}                  -- page 不对，不匹配
	// ✗ {"type":"click"}
	fmt.Println("\n--- 包含查询 (@>) ---")
	rows3, _ := pool.Query(ctx,
		`SELECT id, data->>'type', data->>'page' FROM events 
		 WHERE data @> '{"type":"click","page":"/home"}'`)
	defer rows3.Close()
	for rows3.Next() {
		var id int64
		var t, p string
		rows3.Scan(&id, &t, &p)
		fmt.Printf("  id=%d, type=%s, page=%s\n", id, t, p)
	}

	// ? key 存在判断
	fmt.Println("\n--- key 存在 (?) ---")
	var count int64
	pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM events WHERE data ? 'amount'`,
	).Scan(&count)
	fmt.Printf("  包含 amount 字段的事件数: %d\n", count)

	// 数组字段查询
	fmt.Println("\n--- 数组字段查询 ---")
	rows4, _ := pool.Query(ctx,
		`SELECT id, data->'items' FROM events WHERE data->'items' ? 'book'`)
	defer rows4.Close()
	for rows4.Next() {
		var id int64
		var items string
		rows4.Scan(&id, &items)
		fmt.Printf("  id=%d, items=%s\n", id, items)
	}

	// 聚合统计
	fmt.Println("\n--- 按类型统计 ---")
	rows5, _ := pool.Query(ctx,
		`SELECT data->>'type' AS event_type, COUNT(*) AS cnt 
		 FROM events GROUP BY data->>'type' ORDER BY cnt DESC`)
	defer rows5.Close()
	for rows5.Next() {
		var eventType string
		var cnt int64
		rows5.Scan(&eventType, &cnt)
		fmt.Printf("  %s: %d\n", eventType, cnt)
	}
}

// ============================================================
// 3. JSONB 更新
// ============================================================
func updateOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== JSONB 更新 ==========")

	// || 合并（添加/覆盖字段）把右边的 JSONB 合并进左边。key 不存在就添加，已存在就覆盖值。
	pool.Exec(ctx,
		`UPDATE events SET data = data || '{"processed":true}'::jsonb WHERE id = 1`)
	fmt.Println("id=1 添加 processed:true")

	// - 删除顶层字段，只能删一层，删不了嵌套字段。
	pool.Exec(ctx,
		`UPDATE events SET data = data - 'ip' WHERE data ? 'ip'`)
	fmt.Println("删除所有事件的 ip 字段")

	// #- 删除嵌套字段，按路径删除嵌套字段。{meta,os} 是路径数组，表示删除 data.meta.os。这是 - 的增强版，能深入嵌套结构。
	pool.Exec(ctx,
		`UPDATE events SET data = data #- '{meta,os}' WHERE data ? 'meta'`)
	fmt.Println("删除 meta.os 嵌套字段")

	// jsonb_set 按路径设置值。三个参数：原 JSONB、路径、新值。这里把 data.meta.version 设为 "2.0"。路径不存在时默认会创建。
	pool.Exec(ctx,
		`UPDATE events SET data = jsonb_set(data, '{meta,version}', '"2.0"') WHERE data ? 'meta'`)
	fmt.Println("设置 meta.version = 2.0")

	// 验证
	var data string
	pool.QueryRow(ctx, `SELECT data FROM events WHERE id = 1`).Scan(&data)
	fmt.Printf("id=1 更新后: %s\n", data)
}

// ============================================================
// 4. GIN 索引加速 JSONB 查询
// ============================================================
func ginIndex(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== GIN 索引 ==========")

	// 创建 GIN 索引
	pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_events_data ON events USING GIN (data)`)
	fmt.Println("创建 GIN 索引: idx_events_data")

	// 针对特定字段的表达式索引（更高效）
	pool.Exec(ctx,
		`CREATE INDEX IF NOT EXISTS idx_events_type ON events ((data->>'type'))`)
	fmt.Println("创建表达式索引: idx_events_type")

	// EXPLAIN 验证走索引
	fmt.Println("\nEXPLAIN @> 包含查询:")
	rows, _ := pool.Query(ctx,
		`EXPLAIN SELECT * FROM events WHERE data @> '{"type":"click"}'`)
	defer rows.Close()
	for rows.Next() {
		var plan string
		rows.Scan(&plan)
		fmt.Printf("  %s\n", plan)
	}

	fmt.Println("\nEXPLAIN ->> 等值查询:")
	rows2, _ := pool.Query(ctx,
		`EXPLAIN SELECT * FROM events WHERE data->>'type' = 'click'`)
	defer rows2.Close()
	for rows2.Next() {
		var plan string
		rows2.Scan(&plan)
		fmt.Printf("  %s\n", plan)
	}

	// JSONB 运算符速查
	fmt.Println("\n--- JSONB 运算符速查 ---")
	fmt.Println("  ->    取 JSON 值       data->'meta'")
	fmt.Println("  ->>   取文本值         data->>'type'")
	fmt.Println("  @>    左包含右         data @> '{\"type\":\"click\"}'")
	fmt.Println("  ?     key 存在         data ? 'amount'")
	fmt.Println("  ?|    任一 key 存在    data ?| array['a','b']")
	fmt.Println("  ?&    所有 key 存在    data ?& array['a','b']")
	fmt.Println("  ||    合并             data || '{\"new\":1}'")
	fmt.Println("  -     删除 key         data - 'key'")
	fmt.Println("  #-    删除嵌套 key     data #- '{a,b}'")

	log.Println("[03_jsonb] JSONB 演示完成")
}
