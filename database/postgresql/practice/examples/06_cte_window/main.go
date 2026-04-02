package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/examples/setup"
	pkgPG "pg-practice/pkg/postgres"
)

// 06_cte_window：CTE 递归查询 + 窗口函数
//
// CTE（公用表表达式）：WITH 语法，递归查询树形结构
// 窗口函数：RANK/ROW_NUMBER/SUM OVER 等，不改变行数的聚合
func main() {
	pool := setup.MustSetup()
	defer pkgPG.Close(pool)
	ctx := context.Background()

	cteRecursive(ctx, pool)
	windowFunctions(ctx, pool)
}

// ============================================================
// 1. CTE 递归查询（树形结构）
// ============================================================
func cteRecursive(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("========== CTE 递归查询 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS departments`)
	pool.Exec(ctx, `
		CREATE TABLE departments (
			id        serial PRIMARY KEY,
			name      text NOT NULL,
			parent_id int REFERENCES departments(id)
		)
	`)

	// 插入树形结构
	// 公司
	//   ├── 技术部
	//   │   ├── 后端组
	//   │   │   ├── Go 团队
	//   │   │   └── Java 团队
	//   │   └── 前端组
	//   └── 产品部
	//       └── 设计组
	pool.Exec(ctx, `INSERT INTO departments (id, name, parent_id) VALUES
		(1, '公司', NULL),
		(2, '技术部', 1),
		(3, '产品部', 1),
		(4, '后端组', 2),
		(5, '前端组', 2),
		(6, '设计组', 3),
		(7, 'Go 团队', 4),
		(8, 'Java 团队', 4)
	`)
	fmt.Println("插入部门树（3 层）")

	// 递归查询：查某个部门的所有下级
	fmt.Println("\n--- 技术部的所有下级 ---")
	rows, _ := pool.Query(ctx, `
		WITH RECURSIVE dept_tree AS (
			-- 根节点：技术部
			SELECT id, name, parent_id, 0 AS depth
			FROM departments WHERE id = 2
			
			UNION ALL
			
			-- 递归：查所有子节点
			SELECT d.id, d.name, d.parent_id, dt.depth + 1
			FROM departments d
			JOIN dept_tree dt ON d.parent_id = dt.id
		)
		SELECT id, name, depth, repeat('  ', depth) || name AS display
		FROM dept_tree
		ORDER BY depth, id
	`)
	defer rows.Close()
	for rows.Next() {
		var id, depth int
		var name, display string
		rows.Scan(&id, &name, &depth, &display)
		fmt.Printf("  %s (id=%d, depth=%d)\n", display, id, depth)
	}

	// 递归查询：查某个部门的所有上级（向上追溯）
	fmt.Println("\n--- Go 团队的所有上级（向上追溯）---")
	rows2, _ := pool.Query(ctx, `
		WITH RECURSIVE ancestors AS (
			SELECT id, name, parent_id, 0 AS level
			FROM departments WHERE id = 7
			
			UNION ALL
			
			SELECT d.id, d.name, d.parent_id, a.level + 1
			FROM departments d
			JOIN ancestors a ON d.id = a.parent_id
		)
		SELECT name, level FROM ancestors ORDER BY level
	`)
	defer rows2.Close()
	for rows2.Next() {
		var name string
		var level int
		rows2.Scan(&name, &level)
		fmt.Printf("  level %d: %s\n", level, name)
	}

	// 递归生成序列（数字 1-10）
	fmt.Println("\n--- CTE 生成序列 ---")
	rows3, _ := pool.Query(ctx, `
		WITH RECURSIVE nums AS (
			SELECT 1 AS n
			UNION ALL
			SELECT n + 1 FROM nums WHERE n < 10
		)
		SELECT n FROM nums
	`)
	defer rows3.Close()
	var nums []int
	for rows3.Next() {
		var n int
		rows3.Scan(&n)
		nums = append(nums, n)
	}
	fmt.Printf("  %v\n", nums)
}

// ============================================================
// 2. 窗口函数
// ============================================================
func windowFunctions(ctx context.Context, pool *pgxpool.Pool) {
	fmt.Println("\n========== 窗口函数 ==========")

	pool.Exec(ctx, `DROP TABLE IF EXISTS employees`)
	pool.Exec(ctx, `
		CREATE TABLE employees (
			id     serial PRIMARY KEY,
			name   text NOT NULL,
			dept   text NOT NULL,
			salary bigint NOT NULL
		)
	`)
	pool.Exec(ctx, `INSERT INTO employees (name, dept, salary) VALUES
		('张三', '技术部', 25000),
		('李四', '技术部', 30000),
		('王五', '技术部', 28000),
		('赵六', '产品部', 22000),
		('孙七', '产品部', 26000),
		('周八', '市场部', 20000),
		('吴九', '市场部', 23000),
		('郑十', '市场部', 21000)
	`)

	// RANK：排名（并列有间隔：1,2,2,4）
	fmt.Println("--- 全公司薪资排名（RANK）---")
	rows, _ := pool.Query(ctx, `
		SELECT name, dept, salary,
		       RANK() OVER (ORDER BY salary DESC) AS rk
		FROM employees ORDER BY rk`)
	defer rows.Close()
	for rows.Next() {
		var name, dept string
		var salary int64
		var rk int
		rows.Scan(&name, &dept, &salary, &rk)
		fmt.Printf("  #%d %s (%s): %d\n", rk, name, dept, salary)
	}

	// PARTITION BY：分组内排名
	fmt.Println("\n--- 各部门内薪资排名 ---")
	rows2, _ := pool.Query(ctx, `
		SELECT name, dept, salary,
		       RANK() OVER (PARTITION BY dept ORDER BY salary DESC) AS dept_rank
		FROM employees ORDER BY dept, dept_rank`)
	defer rows2.Close()
	for rows2.Next() {
		var name, dept string
		var salary int64
		var rk int
		rows2.Scan(&name, &dept, &salary, &rk)
		fmt.Printf("  %s #%d %s: %d\n", dept, rk, name, salary)
	}

	// 累计求和
	fmt.Println("\n--- 累计薪资（SUM OVER）---")
	rows3, _ := pool.Query(ctx, `
		SELECT name, salary,
		       SUM(salary) OVER (ORDER BY id) AS running_total
		FROM employees`)
	defer rows3.Close()
	for rows3.Next() {
		var name string
		var salary, total int64
		rows3.Scan(&name, &salary, &total)
		fmt.Printf("  %s: %d, 累计: %d\n", name, salary, total)
	}

	// LAG/LEAD：前后行
	fmt.Println("\n--- 与前一名的薪资差（LAG）---")
	rows4, _ := pool.Query(ctx, `
		SELECT name, salary,
		       LAG(salary, 1) OVER (ORDER BY salary DESC) AS prev_salary,
		       salary - LAG(salary, 1) OVER (ORDER BY salary DESC) AS diff
		FROM employees ORDER BY salary DESC`)
	defer rows4.Close()
	for rows4.Next() {
		var name string
		var salary int64
		var prevSalary, diff *int64
		rows4.Scan(&name, &salary, &prevSalary, &diff)
		d := "—"
		if diff != nil {
			d = fmt.Sprintf("%d", *diff)
		}
		fmt.Printf("  %s: %d, 差距: %s\n", name, salary, d)
	}

	// NTILE：分桶
	fmt.Println("\n--- 薪资分 3 档（NTILE）---")
	rows5, _ := pool.Query(ctx, `
		SELECT name, salary,
		       NTILE(3) OVER (ORDER BY salary DESC) AS tier
		FROM employees ORDER BY tier, salary DESC`)
	defer rows5.Close()
	for rows5.Next() {
		var name string
		var salary int64
		var tier int
		rows5.Scan(&name, &salary, &tier)
		label := map[int]string{1: "高", 2: "中", 3: "低"}
		fmt.Printf("  [%s档] %s: %d\n", label[tier], name, salary)
	}

	log.Println("[06_cte_window] CTE 和窗口函数演示完成")
}
