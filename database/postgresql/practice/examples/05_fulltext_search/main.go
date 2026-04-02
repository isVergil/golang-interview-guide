package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 05_fulltext_search：PG 内置全文检索
//
// 不依赖 Elasticsearch，PG 原生支持 tsvector + tsquery + GIN 索引
// 适合中小规模的搜索场景（百万级数据量完全够用）
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	initData(ctx, pool)
	basicSearch(ctx, pool)
	advancedSearch(ctx, pool)
	ranking(ctx, pool)
}

func initData(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 初始化 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS docs`)
	pool.Exec(ctx, `
		CREATE TABLE docs (
			id      serial PRIMARY KEY,
			title   text NOT NULL,
			content text NOT NULL,
			-- 存储预计算的 tsvector（避免每次查询都重新计算）
			tsv     tsvector GENERATED ALWAYS AS (
				setweight(to_tsvector('english', title), 'A') ||
				setweight(to_tsvector('english', content), 'B')
			) STORED
		)
	`)

	// GIN 索引加速全文检索
	pool.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_docs_tsv ON docs USING GIN (tsv)`)

	pool.Exec(ctx, `INSERT INTO docs (title, content) VALUES
		('Redis Cache Design', 'Redis is an in-memory data store used for caching. It supports various data structures like strings, hashes, and sorted sets.'),
		('PostgreSQL JSONB Guide', 'PostgreSQL provides native JSONB support with GIN indexing. JSONB stores data in binary format for fast queries.'),
		('Go Concurrency Patterns', 'Go provides goroutines and channels for concurrent programming. The select statement allows waiting on multiple channel operations.'),
		('Distributed Lock with Redis', 'Implementing distributed locks using Redis SETNX command with TTL. Lua scripts ensure atomic lock release.'),
		('Database Index Optimization', 'B-Tree indexes are the default in most databases. PostgreSQL also supports GIN, GiST, BRIN, and Hash indexes.'),
		('Microservice Architecture', 'Microservices communicate through REST APIs or gRPC. Service discovery and load balancing are essential components.'),
		('Redis Cluster Sharding', 'Redis Cluster uses 16384 hash slots for data sharding. Each master node handles a subset of slots.'),
		('PostgreSQL MVCC Explained', 'PostgreSQL uses multi-version concurrency control. Each row has xmin and xmax for visibility checking.')
	`)
	fmt.Println("插入 8 篇文档（title 权重 A，content 权重 B）")
}

// ============================================================
// 1. 基础全文检索
// ============================================================
func basicSearch(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 基础全文检索 ==========")

	// 单词搜索
	fmt.Println("--- 搜索 'redis' ---")
	rows, _ := pool.Query(ctx,
		`SELECT id, title FROM docs WHERE tsv @@ to_tsquery('english', 'redis') ORDER BY id`)
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		rows.Scan(&id, &title)
		fmt.Printf("  id=%d: %s\n", id, title)
	}

	// AND 组合搜索
	fmt.Println("\n--- 搜索 'redis & cache' ---")
	rows2, _ := pool.Query(ctx,
		`SELECT id, title FROM docs WHERE tsv @@ to_tsquery('english', 'redis & cache')`)
	defer rows2.Close()
	for rows2.Next() {
		var id int
		var title string
		rows2.Scan(&id, &title)
		fmt.Printf("  id=%d: %s\n", id, title)
	}

	// OR 搜索
	fmt.Println("\n--- 搜索 'redis | postgresql' ---")
	rows3, _ := pool.Query(ctx,
		`SELECT id, title FROM docs WHERE tsv @@ to_tsquery('english', 'redis | postgresql')`)
	defer rows3.Close()
	for rows3.Next() {
		var id int
		var title string
		rows3.Scan(&id, &title)
		fmt.Printf("  id=%d: %s\n", id, title)
	}

	// NOT 排除
	fmt.Println("\n--- 搜索 'index & !redis'（有 index 但没 redis）---")
	rows4, _ := pool.Query(ctx,
		`SELECT id, title FROM docs WHERE tsv @@ to_tsquery('english', 'index & !redis')`)
	defer rows4.Close()
	for rows4.Next() {
		var id int
		var title string
		rows4.Scan(&id, &title)
		fmt.Printf("  id=%d: %s\n", id, title)
	}
}

// ============================================================
// 2. 高级搜索
// ============================================================
func advancedSearch(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 高级搜索 ==========")

	// 短语搜索（相邻词）
	fmt.Println("--- 短语搜索 'distributed <-> lock' ---")
	rows, _ := pool.Query(ctx,
		`SELECT id, title FROM docs WHERE tsv @@ to_tsquery('english', 'distributed <-> lock')`)
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		rows.Scan(&id, &title)
		fmt.Printf("  id=%d: %s\n", id, title)
	}

	// 前缀搜索
	fmt.Println("\n--- 前缀搜索 'micro:*' ---")
	rows2, _ := pool.Query(ctx,
		`SELECT id, title FROM docs WHERE tsv @@ to_tsquery('english', 'micro:*')`)
	defer rows2.Close()
	for rows2.Next() {
		var id int
		var title string
		rows2.Scan(&id, &title)
		fmt.Printf("  id=%d: %s\n", id, title)
	}

	// websearch_to_tsquery：更友好的搜索语法（像 Google 搜索）
	fmt.Println("\n--- websearch 语法 'redis -cluster' ---")
	rows3, _ := pool.Query(ctx,
		`SELECT id, title FROM docs WHERE tsv @@ websearch_to_tsquery('english', 'redis -cluster')`)
	defer rows3.Close()
	for rows3.Next() {
		var id int
		var title string
		rows3.Scan(&id, &title)
		fmt.Printf("  id=%d: %s\n", id, title)
	}

	// 高亮显示匹配片段
	fmt.Println("\n--- 高亮显示 ---")
	rows4, _ := pool.Query(ctx,
		`SELECT id, ts_headline('english', content, to_tsquery('english', 'redis & lock'),
			'StartSel=【, StopSel=】, MaxWords=30') AS highlight
		 FROM docs WHERE tsv @@ to_tsquery('english', 'redis & lock')`)
	defer rows4.Close()
	for rows4.Next() {
		var id int
		var highlight string
		rows4.Scan(&id, &highlight)
		fmt.Printf("  id=%d: %s\n", id, highlight)
	}
}

// ============================================================
// 3. 相关性排名
// ============================================================
func ranking(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 相关性排名 ==========")

	// ts_rank：根据匹配度排序（权重 A > B）
	fmt.Println("--- 搜索 'postgresql' 按相关性排序 ---")
	rows, _ := pool.Query(ctx,
		`SELECT id, title, ts_rank(tsv, to_tsquery('english', 'postgresql')) AS rank
		 FROM docs
		 WHERE tsv @@ to_tsquery('english', 'postgresql')
		 ORDER BY rank DESC`)
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		var rank float64
		rows.Scan(&id, &title, &rank)
		fmt.Printf("  id=%d, rank=%.4f: %s\n", id, rank, title)
	}

	// EXPLAIN 验证走 GIN 索引
	fmt.Println("\n--- EXPLAIN ---")
	rows2, _ := pool.Query(ctx,
		`EXPLAIN SELECT * FROM docs WHERE tsv @@ to_tsquery('english', 'redis')`)
	defer rows2.Close()
	for rows2.Next() {
		var plan string
		rows2.Scan(&plan)
		fmt.Printf("  %s\n", plan)
	}

	log.Println("[05_fulltext_search] 全文检索演示完成")
}
