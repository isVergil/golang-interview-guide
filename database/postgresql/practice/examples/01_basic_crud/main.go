package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 01_basic_crud：pgx 原生 SQL 的增删改查
//
// pgx 是 Go 生态最流行的 PostgreSQL 驱动，纯 Go 实现
// 不使用 ORM，直接写 SQL，性能最优，适合理解 PG 底层操作
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	initTable(ctx, pool)
	insertOps(ctx, pool)
	queryOps(ctx, pool)
	updateOps(ctx, pool)
	deleteOps(ctx, pool)
	batchOps(ctx, pool)
}

func initTable(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 建表 ==========")

	// PG 支持 DDL 事务（CREATE TABLE 也可以回滚，MySQL 不行）
	_, err := pool.Exec(ctx, `
		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			id         bigserial PRIMARY KEY,
			name       text        NOT NULL,
			age        int         NOT NULL DEFAULT 0,
			email      text        NOT NULL UNIQUE,
			balance    bigint      NOT NULL DEFAULT 0,
			created_at timestamptz NOT NULL DEFAULT now(),
			updated_at timestamptz NOT NULL DEFAULT now(),
			deleted_at timestamptz
		)
	`)
	if err != nil {
		log.Fatalf("create table: %v", err)
	}
	fmt.Println("创建 users 表成功")
}

// ============================================================
// 1. INSERT
// ============================================================
func insertOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== INSERT ==========")

	// 单条插入，RETURNING 获取自增 ID（PG 特色，MySQL 需要 LastInsertId）
	var id int64
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, age, email, balance) VALUES ($1, $2, $3, $4) RETURNING id`,
		"alice", 25, "alice@test.com", 10000,
	).Scan(&id)
	if err != nil {
		log.Fatalf("insert: %v", err)
	}
	fmt.Printf("插入 alice, id=%d\n", id)

	// $1, $2 占位符（PG 用 $N，MySQL 用 ?）
	pool.QueryRow(ctx,
		`INSERT INTO users (name, age, email, balance) VALUES ($1, $2, $3, $4) RETURNING id`,
		"bob", 30, "bob@test.com", 5000,
	).Scan(&id)
	fmt.Printf("插入 bob, id=%d\n", id)

	pool.QueryRow(ctx,
		`INSERT INTO users (name, age, email, balance) VALUES ($1, $2, $3, $4) RETURNING id`,
		"charlie", 28, "charlie@test.com", 8000,
	).Scan(&id)
	fmt.Printf("插入 charlie, id=%d\n", id)

	// INSERT ON CONFLICT（PG 版的 Upsert，类似 MySQL 的 ON DUPLICATE KEY UPDATE）
	var upsertID int64
	err = pool.QueryRow(ctx,
		`INSERT INTO users (name, age, email, balance) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name, age = EXCLUDED.age
		 RETURNING id`,
		"alice_v2", 26, "alice@test.com", 10000,
	).Scan(&upsertID)
	if err != nil {
		log.Fatalf("upsert: %v", err)
	}
	fmt.Printf("Upsert alice（email 冲突，更新 name/age）, id=%d\n", upsertID)
}

// ============================================================
// 2. SELECT
// ============================================================
func queryOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== SELECT ==========")

	// 单行查询
	var name string
	var age int
	var balance int64
	err := pool.QueryRow(ctx,
		`SELECT name, age, balance FROM users WHERE email = $1`, "alice@test.com",
	).Scan(&name, &age, &balance)
	if err != nil {
		log.Fatalf("query row: %v", err)
	}
	fmt.Printf("单行查询: name=%s, age=%d, balance=%d\n", name, age, balance)

	// 多行查询
	fmt.Println("\n全部用户:")
	rows, _ := pool.Query(ctx,
		`SELECT id, name, age, email, balance FROM users WHERE deleted_at IS NULL ORDER BY id`)
	defer rows.Close()

	for rows.Next() {
		var id int64
		var n, e string
		var a int
		var b int64
		rows.Scan(&id, &n, &a, &e, &b)
		fmt.Printf("  id=%d, name=%s, age=%d, email=%s, balance=%d\n", id, n, a, e, b)
	}

	// 聚合查询
	var count int64
	var totalBalance int64
	pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(balance), 0) FROM users WHERE deleted_at IS NULL`,
	).Scan(&count, &totalBalance)
	fmt.Printf("\n统计: 总用户=%d, 总余额=%d\n", count, totalBalance)

	// 查看物理行位置 ctid（PG 特色）
	fmt.Println("\n行物理位置 ctid:")
	rows2, _ := pool.Query(ctx, `SELECT ctid, id, name FROM users ORDER BY id`)
	defer rows2.Close()
	for rows2.Next() {
		var ctid string
		var id int64
		var n string
		rows2.Scan(&ctid, &id, &n)
		fmt.Printf("  ctid=%s, id=%d, name=%s\n", ctid, id, n)
	}
}

// ============================================================
// 3. UPDATE
// ============================================================
func updateOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== UPDATE ==========")

	// 普通更新
	tag, _ := pool.Exec(ctx,
		`UPDATE users SET age = $1, updated_at = now() WHERE email = $2`,
		27, "alice@test.com",
	)
	fmt.Printf("更新 alice age=27, 影响行数: %d\n", tag.RowsAffected())

	// 原子更新（余额加减，和 MySQL 的 gorm.Expr 类似）
	tag, _ = pool.Exec(ctx,
		`UPDATE users SET balance = balance + $1, updated_at = now() WHERE email = $2`,
		2000, "bob@test.com",
	)
	fmt.Printf("bob 余额 +2000, 影响行数: %d\n", tag.RowsAffected())

	// RETURNING 更新后的值（PG 特色，不用再查一次）
	var newBalance int64
	pool.QueryRow(ctx,
		`UPDATE users SET balance = balance - $1, updated_at = now() WHERE email = $2 RETURNING balance`,
		1000, "alice@test.com",
	).Scan(&newBalance)
	fmt.Printf("alice 余额 -1000, 更新后余额: %d\n", newBalance)
}

// ============================================================
// 4. DELETE（软删除）
// ============================================================
func deleteOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== DELETE ==========")

	// 软删除（设置 deleted_at）
	tag, _ := pool.Exec(ctx,
		`UPDATE users SET deleted_at = $1 WHERE email = $2`, time.Now(), "charlie@test.com",
	)
	fmt.Printf("软删除 charlie, 影响行数: %d\n", tag.RowsAffected())

	// 验证软删除后查不到
	var count int64
	pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`,
	).Scan(&count)
	fmt.Printf("未删除用户数: %d\n", count)

	// 恢复软删除
	pool.Exec(ctx, `UPDATE users SET deleted_at = NULL WHERE email = $1`, "charlie@test.com")
	fmt.Println("恢复 charlie 软删除")
}

// ============================================================
// 5. 批量操作
// ============================================================
func batchOps(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 批量操作 ==========")

	// 批量插入（unnest 方式，PG 高效批量插入）
	names := []string{"dave", "eve", "frank"}
	ages := []int32{22, 35, 40}
	emails := []string{"dave@test.com", "eve@test.com", "frank@test.com"}
	balances := []int64{3000, 15000, 20000}

	tag, err := pool.Exec(ctx,
		`INSERT INTO users (name, age, email, balance)
		 SELECT * FROM unnest($1::text[], $2::int[], $3::text[], $4::bigint[])`,
		names, ages, emails, balances,
	)
	if err != nil {
		log.Printf("batch insert: %v", err)
	} else {
		fmt.Printf("批量插入 %d 行\n", tag.RowsAffected())
	}

	// 批量更新（用 CTE）
	tag, _ = pool.Exec(ctx,
		`UPDATE users SET balance = balance * 2, updated_at = now() WHERE age >= $1`, 35,
	)
	fmt.Printf("年龄>=35 余额翻倍, 影响行数: %d\n", tag.RowsAffected())

	// 验证
	fmt.Println("\n最终用户列表:")
	rows, _ := pool.Query(ctx,
		`SELECT id, name, age, balance FROM users WHERE deleted_at IS NULL ORDER BY id`)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var n string
		var a int
		var b int64
		rows.Scan(&id, &n, &a, &b)
		fmt.Printf("  id=%d, name=%s, age=%d, balance=%d\n", id, n, a, b)
	}

	log.Println("[01_basic_crud] CRUD 演示完成")
}
