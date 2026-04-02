package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 09_advisory_lock：PostgreSQL 咨询锁（应用层分布式锁）
//
// Advisory Lock 是 PG 提供的应用层锁基础设施：
//   - 不锁表、不锁行，锁的是一个自定义的 bigint key
//   - 业务自己决定 key 的含义（如 task_id、resource_id）
//   - 比 Redis SETNX 的优势：不依赖外部组件、和 PG 事务原子、无 TTL 过期风险
//
// 三种类型：
//   pg_advisory_lock(key)          会话级阻塞锁（手动释放）
//   pg_advisory_xact_lock(key)     事务级锁（事务结束自动释放）
//   pg_try_advisory_lock(key)      非阻塞尝试（返回 bool）
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	initData(ctx, pool)
	sessionLock(ctx, pool)
	txLock(ctx, pool)
	tryLock(ctx, pool)
	concurrentWorkers(ctx, pool)

	log.Println("[09_advisory_lock] 演示完成")
}

func initData(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== 初始化 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS tasks`)
	pool.Exec(ctx, `
		CREATE TABLE tasks (
			id          serial PRIMARY KEY,
			name        text NOT NULL,
			status      text NOT NULL DEFAULT 'pending',
			locked_by   text,
			executed_at timestamptz
		)
	`)

	pool.Exec(ctx, `
		INSERT INTO tasks (name) VALUES
		('发送日报邮件'),
		('同步用户数据'),
		('清理过期缓存'),
		('生成统计报表')
	`)
	fmt.Println("创建 tasks 表并插入 4 个任务")
}

// ============================================================
// 1. 会话级锁（手动获取和释放）
// ============================================================
func sessionLock(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 会话级锁 ==========")
	fmt.Println("pg_advisory_lock(key): 阻塞直到获取锁")
	fmt.Println("pg_advisory_unlock(key): 手动释放")

	lockKey := int64(1001)

	// 获取锁
	_, err := pool.Exec(ctx, `SELECT pg_advisory_lock($1)`, lockKey)
	if err != nil {
		log.Printf("获取锁失败: %v", err)
		return
	}
	fmt.Printf("获取会话锁 key=%d 成功\n", lockKey)

	// 模拟业务操作
	fmt.Println("执行业务操作...")
	time.Sleep(100 * time.Millisecond)

	// 释放锁
	pool.Exec(ctx, `SELECT pg_advisory_unlock($1)`, lockKey)
	fmt.Printf("释放会话锁 key=%d\n", lockKey)

	// 查看当前锁状态
	var count int64
	pool.QueryRow(ctx, `SELECT count(*) FROM pg_locks WHERE locktype = 'advisory'`).Scan(&count)
	fmt.Printf("当前 advisory lock 数量: %d\n", count)
}

// ============================================================
// 2. 事务级锁（事务结束自动释放）
// ============================================================
func txLock(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 事务级锁 ==========")
	fmt.Println("pg_advisory_xact_lock(key): 事务提交/回滚时自动释放，无需手动 unlock")

	lockKey := int64(2001)

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Printf("begin tx: %v", err)
		return
	}

	// 在事务中获取锁
	tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, lockKey)
	fmt.Printf("事务内获取锁 key=%d\n", lockKey)

	// 业务操作
	tx.Exec(ctx, `UPDATE tasks SET status = 'running', executed_at = now() WHERE id = 1`)
	fmt.Println("更新 task 1 状态为 running")

	tx.Exec(ctx, `UPDATE tasks SET status = 'done' WHERE id = 1`)
	fmt.Println("更新 task 1 状态为 done")

	// 提交事务，锁自动释放
	tx.Commit(ctx)
	fmt.Println("事务提交，锁自动释放")
}

// ============================================================
// 3. 非阻塞尝试锁
// ============================================================
func tryLock(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 非阻塞尝试锁 ==========")
	fmt.Println("pg_try_advisory_lock(key): 立即返回 bool，不阻塞等待")

	lockKey := int64(3001)

	// 先获取锁
	pool.Exec(ctx, `SELECT pg_advisory_lock($1)`, lockKey)
	fmt.Printf("先获取锁 key=%d\n", lockKey)

	// 另一个连接尝试获取同一把锁
	var acquired bool
	pool.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, lockKey).Scan(&acquired)
	fmt.Printf("再次尝试获取锁 key=%d: acquired=%v（预期 false，因为锁已被持有）\n", lockKey, acquired)

	// 释放锁
	pool.Exec(ctx, `SELECT pg_advisory_unlock($1)`, lockKey)
	fmt.Printf("释放锁 key=%d\n", lockKey)

	// 现在再尝试
	pool.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, lockKey).Scan(&acquired)
	fmt.Printf("释放后再尝试获取: acquired=%v（预期 true）\n", acquired)
	pool.Exec(ctx, `SELECT pg_advisory_unlock($1)`, lockKey)
}

// ============================================================
// 4. 并发任务调度（实战场景）
// ============================================================
func concurrentWorkers(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 并发任务调度 ==========")
	fmt.Println("场景：多个 worker 竞争执行任务，用 advisory lock 保证同一任务不被重复执行")

	// 重置任务状态
	pool.Exec(ctx, `UPDATE tasks SET status = 'pending', locked_by = NULL, executed_at = NULL`)

	var wg sync.WaitGroup
	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			processTask(ctx, pool, workerID)
		}(i)
	}
	wg.Wait()

	// 查看最终结果
	fmt.Println("\n--- 任务执行结果 ---")
	rows, _ := pool.Query(ctx,
		`SELECT id, name, status, COALESCE(locked_by, '-'), executed_at FROM tasks ORDER BY id`)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name, status, lockedBy string
		var executedAt *time.Time
		rows.Scan(&id, &name, &status, &lockedBy, &executedAt)
		t := "-"
		if executedAt != nil {
			t = executedAt.Format("15:04:05.000")
		}
		fmt.Printf("  task=%d [%s] %s | worker=%s | time=%s\n", id, status, name, lockedBy, t)
	}
}

func processTask(ctx context.Context, pool *pgxpool.Pool, workerID int) {
	// 查询待执行的任务
	rows, _ := pool.Query(ctx, `SELECT id, name FROM tasks WHERE status = 'pending' ORDER BY id`)
	defer rows.Close()

	for rows.Next() {
		var taskID int
		var taskName string
		rows.Scan(&taskID, &taskName)

		// 尝试获取该任务的锁（用 taskID 作为 lock key）
		var acquired bool
		pool.QueryRow(ctx,
			`SELECT pg_try_advisory_lock($1)`, int64(taskID+10000),
		).Scan(&acquired)

		if !acquired {
			// 锁被其他 worker 持有，跳过
			continue
		}

		// 获得锁，执行任务
		fmt.Printf("  worker-%d 获取 task=%d (%s)\n", workerID, taskID, taskName)

		pool.Exec(ctx,
			`UPDATE tasks SET status = 'done', locked_by = $1, executed_at = now() WHERE id = $2 AND status = 'pending'`,
			fmt.Sprintf("worker-%d", workerID), taskID,
		)

		// 模拟执行
		time.Sleep(50 * time.Millisecond)

		// 释放锁
		pool.Exec(ctx, `SELECT pg_advisory_unlock($1)`, int64(taskID+10000))
	}
}
