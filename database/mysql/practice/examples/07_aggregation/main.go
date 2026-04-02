package main

import (
	"context"
	"fmt"

	"mysql-practice/examples/setup"
	"mysql-practice/internal/model"
	"mysql-practice/pkg/database"
)

// 07_aggregation：聚合查询 + Raw SQL
//
// GORM 的 Group/Having/Select 做聚合统计
// 以及当 GORM 表达力不够时，直接用 Raw SQL
func main() {
	db := setup.MustSetup()
	defer database.Close(db)

	ctx := context.Background()

	// 准备测试数据
	fmt.Println("===== 准备数据 =====")
	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引导致后续插入冲突
	db.Unscoped().Where("1 = 1").Delete(&model.User{})
	users := []model.User{
		{Name: "张三", Age: 22, Email: "zs@test.com", Balance: 1000},
		{Name: "李四", Age: 22, Email: "ls@test.com", Balance: 3000},
		{Name: "王五", Age: 28, Email: "ww@test.com", Balance: 15000},
		{Name: "赵六", Age: 28, Email: "zl@test.com", Balance: 8000},
		{Name: "孙七", Age: 35, Email: "sq@test.com", Balance: 20000},
		{Name: "周八", Age: 35, Email: "zb@test.com", Balance: 50000},
		{Name: "吴九", Age: 35, Email: "wj@test.com", Balance: 12000},
	}
	db.Create(&users)
	fmt.Printf("插入 %d 条数据\n", len(users))

	// ---------------------------------------------------------------
	// 一、GORM 聚合：按年龄分组统计
	// SQL: SELECT age, COUNT(*) as cnt, SUM(balance) as total_balance,
	//             AVG(balance) as avg_balance
	//      FROM users GROUP BY age
	// ---------------------------------------------------------------
	fmt.Println("\n===== GORM 聚合：按年龄分组统计 =====")

	type AgeStats struct {
		Age          int   `gorm:"column:age"`
		Count        int   `gorm:"column:cnt"`
		TotalBalance int64 `gorm:"column:total_balance"`
		AvgBalance   int64 `gorm:"column:avg_balance"`
	}

	var stats []AgeStats
	db.WithContext(ctx).Model(&model.User{}).
		Select("age, COUNT(*) as cnt, SUM(balance) as total_balance, AVG(balance) as avg_balance").
		Group("age").
		Order("age").
		Scan(&stats)

	fmt.Printf("%-6s %-6s %-14s %-14s\n", "年龄", "人数", "总余额(分)", "平均余额(分)")
	for _, s := range stats {
		fmt.Printf("%-6d %-6d %-14d %-14d\n", s.Age, s.Count, s.TotalBalance, s.AvgBalance)
	}

	// ---------------------------------------------------------------
	// 二、GORM Having：过滤分组结果
	// SQL: SELECT age, COUNT(*) as cnt FROM users
	//      GROUP BY age HAVING cnt >= 2
	// ---------------------------------------------------------------
	fmt.Println("\n===== Having：人数 >= 2 的年龄段 =====")

	type AgeCnt struct {
		Age   int `gorm:"column:age"`
		Count int `gorm:"column:cnt"`
	}

	var filtered []AgeCnt
	db.WithContext(ctx).Model(&model.User{}).
		Select("age, COUNT(*) as cnt").
		Group("age").
		Having("cnt >= ?", 2).
		Scan(&filtered)

	for _, s := range filtered {
		fmt.Printf("  age=%d, count=%d\n", s.Age, s.Count)
	}

	// ---------------------------------------------------------------
	// 三、Raw SQL：复杂查询直接写 SQL
	// 场景：报表、多表联查、窗口函数等 GORM 表达不了的情况
	// ---------------------------------------------------------------
	fmt.Println("\n===== Raw SQL：余额排名 =====")

	type RankResult struct {
		Name    string `gorm:"column:name"`
		Balance int64  `gorm:"column:balance"`
		Rank    int    `gorm:"column:rk"`
	}

	var ranks []RankResult
	db.WithContext(ctx).Raw(`
		SELECT name, balance,
		       RANK() OVER (ORDER BY balance DESC) as rk
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY rk
	`).Scan(&ranks)

	fmt.Printf("%-8s %-12s %-6s\n", "姓名", "余额(分)", "排名")
	for _, r := range ranks {
		fmt.Printf("%-8s %-12d %-6d\n", r.Name, r.Balance, r.Rank)
	}

	// ---------------------------------------------------------------
	// 四、Raw SQL 执行写操作（谨慎使用）
	// ---------------------------------------------------------------
	fmt.Println("\n===== Raw Exec：批量调整余额 =====")

	result := db.WithContext(ctx).Exec(
		"UPDATE users SET balance = balance * 2 WHERE age >= ? AND deleted_at IS NULL", 35,
	)
	fmt.Printf("影响行数: %d\n", result.RowsAffected)

	// 验证
	var updated []model.User
	db.WithContext(ctx).Where("age >= ?", 35).Find(&updated)
	for _, u := range updated {
		fmt.Printf("  %s: balance=%d (翻倍后)\n", u.Name, u.Balance)
	}

	// ---------------------------------------------------------------
	// 五、Select 指定列（大表优化）
	// ---------------------------------------------------------------
	fmt.Println("\n===== Select 指定列 =====")

	type UserBrief struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}

	var briefs []UserBrief
	db.WithContext(ctx).Model(&model.User{}).
		Select("id", "name").
		Find(&briefs)

	for _, b := range briefs {
		fmt.Printf("  id=%d, name=%s\n", b.ID, b.Name)
	}
	fmt.Println("→ 只查需要的列，减少网络传输和内存占用")

	fmt.Println("\n===== 聚合查询演示完成 =====")
}
