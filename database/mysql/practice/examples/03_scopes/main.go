// 03_scopes：GORM Scopes 可复用查询条件
//
// Scope 是一个 func(db *gorm.DB) *gorm.DB 函数
// 可以链式组合，实现查询条件的复用和标准化
//
// 运行：cd practice && go run examples/03_scopes/main.go
package main

import (
	"context"
	"fmt"

	"mysql-practice/examples/setup"
	"mysql-practice/internal/dao"
	"mysql-practice/internal/model"
	"mysql-practice/pkg/database"
)

func main() {
	db := setup.MustSetup()
	defer database.Close(db)

	userDAO := dao.NewUserDAO(db)
	ctx := context.Background()

	// 准备测试数据
	fmt.Println("===== 准备测试数据 =====")
	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引导致后续插入冲突
	// 业务代码中的删除操作仍然使用软删除
	db.Unscoped().Where("1 = 1").Delete(&model.User{})
	users := []model.User{
		{Name: "张三", Age: 22, Email: "zhangsan@test.com", Balance: 1000},
		{Name: "李四", Age: 28, Email: "lisi@test.com", Balance: 5000},
		{Name: "王五", Age: 35, Email: "wangwu@test.com", Balance: 15000},
		{Name: "赵六", Age: 25, Email: "zhaoliu@test.com", Balance: 800},
		{Name: "孙七", Age: 30, Email: "sunqi@test.com", Balance: 20000},
	}
	userDAO.BatchCreate(ctx, users, 100)
	fmt.Printf("插入 %d 条测试数据\n", len(users))

	// ---------------------------------------------------------------
	// Scope 组合查询演示
	// ---------------------------------------------------------------

	// 查询 1：年龄 25~30 的用户
	fmt.Println("\n===== 年龄 25~30 =====")
	result, _ := userDAO.ListWithScopes(ctx,
		dao.AgeRange(25, 30),
		dao.OrderByCreated(true),
	)
	for _, u := range result {
		fmt.Printf("  %s, age=%d, balance=%d\n", u.Name, u.Age, u.Balance)
	}

	// 查询 2：余额 >= 5000 + 分页第1页（每页2条）
	fmt.Println("\n===== 余额 >= 5000，分页(1, 2) =====")
	result, _ = userDAO.ListWithScopes(ctx,
		dao.BalanceGte(5000),
		dao.Paginate(1, 2),
		dao.OrderByCreated(false),
	)
	for _, u := range result {
		fmt.Printf("  %s, age=%d, balance=%d\n", u.Name, u.Age, u.Balance)
	}

	// 查询 3：组合三个条件
	fmt.Println("\n===== 年龄 20~35 + 余额 >= 1000 + 分页(1, 10) =====")
	result, _ = userDAO.ListWithScopes(ctx,
		dao.AgeRange(20, 35),
		dao.BalanceGte(1000),
		dao.Paginate(1, 10),
	)
	for _, u := range result {
		fmt.Printf("  %s, age=%d, balance=%d\n", u.Name, u.Age, u.Balance)
	}

	// ---------------------------------------------------------------
	// Paginate 的参数保护
	// ---------------------------------------------------------------
	fmt.Println("\n===== Paginate 参数保护 =====")
	fmt.Println("Paginate(0, -1)  → 自动修正为 page=1, size=10")
	fmt.Println("Paginate(1, 999) → 自动修正为 size=100（上限保护）")
	result, _ = userDAO.ListWithScopes(ctx, dao.Paginate(0, -1))
	fmt.Printf("实际返回 %d 条\n", len(result))

	fmt.Println("\n===== Scopes 演示完成 =====")
}
