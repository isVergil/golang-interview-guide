package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 02_transaction：事务隔离级别 + 悲观锁 + 转账
//
// PG 事务特点：
//   - DDL 也能回滚（CREATE TABLE 放在事务里可以回滚，MySQL 不行）
//   - Repeatable Read 没有幻读（真正的快照隔离）
//   - 支持 Serializable（SSI 可序列化快照隔离）
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	initData(ctx, pool)
	ddlTransaction(ctx, pool)
	normalTransfer(ctx, pool)
	failTransfer(ctx, pool)
	concurrentTransfer(ctx, pool)
}

func initData(ctx context.Context, pool *pgxpool.Pool) {
	pool.Exec(ctx, `DROP TABLE IF EXISTS accounts`)
	pool.Exec(ctx, `
		CREATE TABLE accounts (
			id      bigserial PRIMARY KEY,
			name    text   NOT NULL UNIQUE,
			balance bigint NOT NULL DEFAULT 0
		)
	`)
	pool.Exec(ctx, `INSERT INTO accounts (name, balance) VALUES ('alice', 10000), ('bob', 5000)`)
	fmt.Println("========== 初始化 ==========")
	fmt.Println("alice: 10000, bob: 5000")
}

// ============================================================
// 1. DDL 事务（PG 独有）
// ============================================================
func ddlTransaction(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== DDL 事务（PG 独有）==========")

	// 在事务中 CREATE TABLE 然后回滚 → 表不会被创建
	tx, _ := pool.Begin(ctx)

	tx.Exec(ctx, `CREATE TABLE temp_test (id serial PRIMARY KEY, name text)`)
	tx.Exec(ctx, `INSERT INTO temp_test (name) VALUES ('hello')`)
	fmt.Println("事务内: CREATE TABLE + INSERT 完成")

	// 回滚 → 表和数据都不存在
	tx.Rollback(ctx)
	fmt.Println("ROLLBACK → 表 temp_test 不会被创建")

	// 验证
	var exists bool
	pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'temp_test')`,
	).Scan(&exists)
	fmt.Printf("temp_test 是否存在: %v (预期 false)\n", exists)
}

// ============================================================
// 2. 正常转账（FOR UPDATE 悲观锁）
// ============================================================
func normalTransfer(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 正常转账 ==========")

	err := transfer(ctx, pool, "alice", "bob", 3000)
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
		return
	}

	printBalances(ctx, pool)
}

// ============================================================
// 3. 余额不足 → 自动回滚
// ============================================================
func failTransfer(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 余额不足 ==========")

	err := transfer(ctx, pool, "alice", "bob", 99999)
	fmt.Printf("转账结果: %v (预期失败)\n", err)

	printBalances(ctx, pool)
	fmt.Println("→ 余额不变，已自动回滚")
}

// ============================================================
// 4. 并发转账（FOR UPDATE 保证一致性）
// ============================================================
func concurrentTransfer(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 并发转账 ==========")

	// 重置余额
	pool.Exec(ctx, `UPDATE accounts SET balance = 10000 WHERE name = 'alice'`)
	pool.Exec(ctx, `UPDATE accounts SET balance = 10000 WHERE name = 'bob'`)
	fmt.Println("重置: alice=10000, bob=10000")

	var wg sync.WaitGroup
	var successCount int64
	var failCount int64

	// 10 个 goroutine 并发 alice → bob 转 1000
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := transfer(ctx, pool, "alice", "bob", 1000)
			if err != nil {
				atomic.AddInt64(&failCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}()
	}
	wg.Wait()

	fmt.Printf("10 次并发转账: 成功=%d, 失败=%d\n", successCount, failCount)
	printBalances(ctx, pool)

	var aliceBalance, bobBalance int64
	pool.QueryRow(ctx, `SELECT balance FROM accounts WHERE name='alice'`).Scan(&aliceBalance)
	pool.QueryRow(ctx, `SELECT balance FROM accounts WHERE name='bob'`).Scan(&bobBalance)
	fmt.Printf("总和: %d (应为 20000，资金守恒)\n", aliceBalance+bobBalance)

	log.Println("[02_transaction] 事务演示完成")
}

// transfer 转账（事务 + FOR UPDATE 悲观锁）
func transfer(ctx context.Context, pool *pgxpool.Pool, from, to string, amount int64) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// FOR UPDATE 锁住转出方（按 name 排序加锁，防死锁）
	var fromBalance int64
	err = tx.QueryRow(ctx,
		`SELECT balance FROM accounts WHERE name = $1 FOR UPDATE`, from,
	).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("lock from: %w", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("余额不足: %d < %d", fromBalance, amount)
	}

	// FOR UPDATE 锁住转入方
	tx.QueryRow(ctx,
		`SELECT balance FROM accounts WHERE name = $1 FOR UPDATE`, to,
	)

	// 原子扣减/增加
	tx.Exec(ctx, `UPDATE accounts SET balance = balance - $1 WHERE name = $2`, amount, from)
	tx.Exec(ctx, `UPDATE accounts SET balance = balance + $1 WHERE name = $2`, amount, to)

	return tx.Commit(ctx)
}

func printBalances(ctx context.Context, pool *pgxpool.Pool) {
	rows, _ := pool.Query(ctx, `SELECT name, balance FROM accounts ORDER BY name`)
	defer rows.Close()
	for rows.Next() {
		var name string
		var balance int64
		rows.Scan(&name, &balance)
		fmt.Printf("  %s: %d\n", name, balance)
	}
}
