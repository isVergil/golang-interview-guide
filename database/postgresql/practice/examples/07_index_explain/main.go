package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 07_index_explain：索引类型与 EXPLAIN 分析
//
// PostgreSQL 支持 6 种索引（B-Tree、Hash、GIN、GiST、BRIN、SP-GiST）
// 远多于 MySQL（基本只有 B+Tree）
// 本例演示各类索引的创建和 EXPLAIN 查询计划分析
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	initData(ctx, pool)
	btreeIndex(ctx, pool)
	hashIndex(ctx, pool)
	ginIndex(ctx, pool)
	brinIndex(ctx, pool)
	partialAndCoveringIndex(ctx, pool)
	explainDemo(ctx, pool)

	log.Println("[07_index_explain] 演示完成")
}

func initData(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 初始化数据 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS products`)
	pool.Exec(ctx, `DROP TABLE IF EXISTS logs`)
	pool.Exec(ctx, `
		CREATE TABLE products (
			id       bigserial PRIMARY KEY,
			name     text NOT NULL,
			category text NOT NULL,
			price    numeric(10,2) NOT NULL,
			tags     text[],
			meta     jsonb,
			status   smallint NOT NULL DEFAULT 1,
			created_at timestamptz NOT NULL DEFAULT now()
		)
	`)

	// 插入测试数据
	pool.Exec(ctx, `
		INSERT INTO products (name, category, price, tags, meta, status)
		SELECT
			'product_' || i,
			CASE WHEN i % 3 = 0 THEN 'electronics'
			     WHEN i % 3 = 1 THEN 'books'
			     ELSE 'clothing' END,
			(random() * 1000)::numeric(10,2),
			CASE WHEN i % 2 = 0 THEN ARRAY['hot','new']
			     ELSE ARRAY['sale'] END,
			jsonb_build_object('weight', (random()*10)::numeric(4,1), 'color', 
				CASE WHEN i%3=0 THEN 'red' WHEN i%3=1 THEN 'blue' ELSE 'green' END),
			CASE WHEN i % 10 = 0 THEN 0 ELSE 1 END
		FROM generate_series(1, 10000) AS i
	`)
	fmt.Println("插入 10000 条 products 数据")

	// 时序表（用于 BRIN 索引）
	pool.Exec(ctx, `
		CREATE TABLE logs (
			id         bigserial PRIMARY KEY,
			level      text NOT NULL,
			message    text NOT NULL,
			created_at timestamptz NOT NULL DEFAULT now()
		)
	`)
	pool.Exec(ctx, `
		INSERT INTO logs (level, message, created_at)
		SELECT
			CASE WHEN i%5=0 THEN 'ERROR' WHEN i%3=0 THEN 'WARN' ELSE 'INFO' END,
			'log message ' || i,
			now() - ((10000-i) || ' seconds')::interval
		FROM generate_series(1, 10000) AS i
	`)
	fmt.Println("插入 10000 条 logs 数据（按时间顺序）")

	// 更新统计信息，让优化器有准确的数据
	pool.Exec(ctx, `ANALYZE products`)
	pool.Exec(ctx, `ANALYZE logs`)
}

// ============================================================
// 1. B-Tree 索引（默认）
// ============================================================
func btreeIndex(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== B-Tree 索引 ==========")
	fmt.Println("适合：=, <, >, BETWEEN, ORDER BY, IS NULL")

	// 普通 B-Tree
	pool.Exec(ctx, `CREATE INDEX idx_products_category ON products (category)`)
	fmt.Println("创建 B-Tree 索引: idx_products_category")

	// 复合索引（最左前缀原则）
	pool.Exec(ctx, `CREATE INDEX idx_products_cat_price ON products (category, price)`)
	fmt.Println("创建复合索引: idx_products_cat_price (category, price)")

	// 验证最左前缀
	fmt.Println("\n--- 最左前缀验证 ---")
	printExplain(ctx, pool, "走索引（最左列）",
		`SELECT * FROM products WHERE category = 'books'`)
	printExplain(ctx, pool, "走索引（最左 + 后续列）",
		`SELECT * FROM products WHERE category = 'books' AND price > 500`)
	printExplain(ctx, pool, "不走复合索引（跳过最左列）",
		`SELECT * FROM products WHERE price > 500`)
}

// ============================================================
// 2. Hash 索引
// ============================================================
func hashIndex(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== Hash 索引 ==========")
	fmt.Println("适合：纯等值查询 (=)，不支持范围、排序")

	pool.Exec(ctx, `CREATE INDEX idx_products_name_hash ON products USING HASH (name)`)
	fmt.Println("创建 Hash 索引: idx_products_name_hash")

	printExplain(ctx, pool, "Hash 索引等值查询",
		`SELECT * FROM products WHERE name = 'product_100'`)

	// 查看索引大小对比
	fmt.Println("\n--- 索引大小对比 ---")
	rows, _ := pool.Query(ctx, `
		SELECT indexname, pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
		FROM pg_indexes WHERE tablename = 'products' ORDER BY indexname`)
	defer rows.Close()
	for rows.Next() {
		var name, size string
		rows.Scan(&name, &size)
		fmt.Printf("  %s: %s\n", name, size)
	}
}

// ============================================================
// 3. GIN 索引（倒排索引）
// ============================================================
func ginIndex(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== GIN 索引 ==========")
	fmt.Println("适合：JSONB (@>, ?), 数组 (@>, &&), 全文检索 (@@)")

	// JSONB GIN 索引
	pool.Exec(ctx, `CREATE INDEX idx_products_meta ON products USING GIN (meta)`)
	fmt.Println("创建 GIN 索引: idx_products_meta (JSONB)")

	// 数组 GIN 索引
	pool.Exec(ctx, `CREATE INDEX idx_products_tags ON products USING GIN (tags)`)
	fmt.Println("创建 GIN 索引: idx_products_tags (数组)")

	// JSONB 包含查询
	printExplain(ctx, pool, "GIN 索引 JSONB 包含查询",
		`SELECT count(*) FROM products WHERE meta @> '{"color":"red"}'`)

	// 数组包含查询
	printExplain(ctx, pool, "GIN 索引数组包含查询",
		`SELECT count(*) FROM products WHERE tags @> ARRAY['hot']`)
}

// ============================================================
// 4. BRIN 索引（块范围索引）
// ============================================================
func brinIndex(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== BRIN 索引 ==========")
	fmt.Println("适合：数据物理有序的大表（时序数据）")
	fmt.Println("原理：不索引每一行，只记录每个页范围的 min/max，体积极小")

	// BRIN 索引
	pool.Exec(ctx, `CREATE INDEX idx_logs_time_brin ON logs USING BRIN (created_at)`)
	fmt.Println("创建 BRIN 索引: idx_logs_time_brin")

	// 同列创建 B-Tree 做对比
	pool.Exec(ctx, `CREATE INDEX idx_logs_time_btree ON logs (created_at)`)
	fmt.Println("创建 B-Tree 索引: idx_logs_time_btree (对比)")

	// 大小对比
	fmt.Println("\n--- BRIN vs B-Tree 大小对比 ---")
	rows, _ := pool.Query(ctx, `
		SELECT indexname, pg_size_pretty(pg_relation_size(indexname::regclass)) AS size
		FROM pg_indexes WHERE tablename = 'logs' ORDER BY indexname`)
	defer rows.Close()
	for rows.Next() {
		var name, size string
		rows.Scan(&name, &size)
		fmt.Printf("  %s: %s\n", name, size)
	}
	fmt.Println("  → BRIN 远小于 B-Tree（万行级就有明显差距，亿行级差 2 万倍）")

	// 范围查询
	printExplain(ctx, pool, "BRIN 索引时间范围查询",
		`SELECT count(*) FROM logs WHERE created_at > now() - interval '1 hour'`)
}

// ============================================================
// 5. 部分索引 & 覆盖索引
// ============================================================
func partialAndCoveringIndex(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 部分索引 & 覆盖索引 ==========")

	// 部分索引：只索引满足条件的行，减小体积
	pool.Exec(ctx, `CREATE INDEX idx_products_active ON products (name) WHERE status = 1`)
	fmt.Println("创建部分索引: idx_products_active (仅 status=1 的行)")

	printExplain(ctx, pool, "部分索引查询 (status=1)",
		`SELECT name FROM products WHERE status = 1 AND name = 'product_50'`)

	// 覆盖索引：INCLUDE 携带额外列，避免回堆表
	pool.Exec(ctx, `CREATE INDEX idx_products_cat_cover ON products (category) INCLUDE (name, price)`)
	fmt.Println("\n创建覆盖索引: idx_products_cat_cover (category) INCLUDE (name, price)")

	printExplain(ctx, pool, "覆盖索引 Index Only Scan",
		`SELECT name, price FROM products WHERE category = 'books'`)

	// 表达式索引
	pool.Exec(ctx, `CREATE INDEX idx_products_lower_name ON products (lower(name))`)
	fmt.Println("\n创建表达式索引: idx_products_lower_name (lower(name))")

	printExplain(ctx, pool, "表达式索引",
		`SELECT * FROM products WHERE lower(name) = 'product_100'`)
}

// ============================================================
// 6. EXPLAIN 详解
// ============================================================
func explainDemo(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== EXPLAIN 详解 ==========")
	fmt.Println(`
EXPLAIN 关键字段：
  cost=启动代价..总代价    估算值（单位：page fetch 开销）
  rows=N                 估算返回行数
  actual time=ms..ms     实际执行时间
  Buffers: shared hit=N  缓存命中页数
  Buffers: shared read=N 磁盘读取页数

常见扫描方式：
  Seq Scan           全表扫描（小表正常，大表需优化）
  Index Scan         索引扫描 → 回堆表取数据
  Index Only Scan    覆盖索引，不回表（最快）
  Bitmap Index Scan  索引构建位图 → 批量回表（多条件 OR / 返回行多时）`)

	// 全表扫描 vs 索引扫描
	fmt.Println("\n--- Seq Scan vs Index Scan ---")
	printExplainVerbose(ctx, pool,
		`SELECT * FROM products WHERE category = 'books' AND price > 800`)

	// Bitmap Scan（返回行较多时）
	fmt.Println("\n--- Bitmap Scan（返回行较多时优化器自动选择）---")
	printExplainVerbose(ctx, pool,
		`SELECT * FROM products WHERE tags @> ARRAY['hot'] OR tags @> ARRAY['sale']`)
}

// ============================================================
// 工具函数
// ============================================================

// printExplain 打印简化的 EXPLAIN 输出
func printExplain(ctx context.Context, pool *pgxpool.Pool, label, query string) {
	rows, err := pool.Query(ctx, "EXPLAIN "+query)
	if err != nil {
		fmt.Printf("  [%s] explain error: %v\n", label, err)
		return
	}
	defer rows.Close()

	fmt.Printf("  [%s]\n", label)
	for rows.Next() {
		var line string
		rows.Scan(&line)
		fmt.Printf("    %s\n", line)
	}
}

// printExplainVerbose 打印 EXPLAIN (ANALYZE, BUFFERS) 的完整输出
func printExplainVerbose(ctx context.Context, pool *pgxpool.Pool, query string) {
	rows, err := pool.Query(ctx, "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) "+query)
	if err != nil {
		fmt.Printf("  explain error: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("  SQL: %s\n", strings.TrimSpace(query))
	for rows.Next() {
		var line string
		rows.Scan(&line)
		fmt.Printf("    %s\n", line)
	}
}
