package main

import (
	"context"
	"fmt"

	"mysql-practice/examples/setup"
	"mysql-practice/internal/dao"
	"mysql-practice/internal/model"
	"mysql-practice/pkg/database"
	"mysql-practice/pkg/snowflake"
)

// 04_association：关联查询 Preload vs Joins
//
// Preload：两次查询，先查主表再查关联表（适合一对多）
// Joins：单次 JOIN 查询（适合一对一 / 需要关联字段做 WHERE）
func main() {
	db := setup.MustSetup()
	defer database.Close(db)

	userDAO := dao.NewUserDAO(db)
	orderDAO := dao.NewOrderDAO(db)
	ctx := context.Background()

	// 初始化雪花节点
	node, _ := snowflake.NewNode(1)

	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引和外键冲突
	// 先删 orders（有外键引用 users），再删 users
	db.Unscoped().Where("1 = 1").Delete(&model.Order{})
	db.Unscoped().Where("1 = 1").Delete(&model.User{})

	// 准备用户
	fmt.Println("===== 准备数据 =====")
	alice := &model.User{Name: "alice", Age: 25, Email: "alice_assoc@test.com"}
	bob := &model.User{Name: "bob", Age: 30, Email: "bob_assoc@test.com"}
	userDAO.Create(ctx, alice)
	userDAO.Create(ctx, bob)
	fmt.Printf("用户: alice(id=%d), bob(id=%d)\n", alice.ID, bob.ID)

	// 准备订单（雪花 ID 作为业务订单号，自增 ID 作为主键）
	orders := []model.Order{
		{OrderNo: node.Generate(), UserID: alice.ID, ProductName: "Go 语言圣经", Amount: 6800, Status: 1},
		{OrderNo: node.Generate(), UserID: alice.ID, ProductName: "Redis 设计与实现", Amount: 5500, Status: 0},
		{OrderNo: node.Generate(), UserID: bob.ID, ProductName: "MySQL 技术内幕", Amount: 7900, Status: 1},
	}
	for i := range orders {
		orderDAO.Create(ctx, &orders[i])
	}
	fmt.Printf("创建 %d 个订单\n", len(orders))

	// ---------------------------------------------------------------
	// Preload：按业务订单号查询，自动加载关联的 User
	// SQL: SELECT * FROM orders WHERE order_no = ?;
	//      SELECT * FROM users WHERE id IN (?);
	// ---------------------------------------------------------------
	fmt.Println("\n===== Preload：按订单号查订单 =====")
	order, err := orderDAO.GetByOrderNo(ctx, orders[0].OrderNo)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("订单号: %d (主键ID: %d)\n", order.OrderNo, order.ID)
		fmt.Printf("商品: %s, 金额: %d分\n", order.ProductName, order.Amount)
		fmt.Printf("用户: %s, Email: %s\n", order.User.Name, order.User.Email)
	}

	// ---------------------------------------------------------------
	// Preload：查用户的所有订单
	// ---------------------------------------------------------------
	fmt.Println("\n===== Preload：查 alice 的所有订单 =====")
	aliceOrders, err := orderDAO.ListByUserID(ctx, alice.ID)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		for _, o := range aliceOrders {
			fmt.Printf("  订单号: %d, 商品: %s, 金额: %d分, 状态: %d\n",
				o.OrderNo, o.ProductName, o.Amount, o.Status)
		}
	}

	// ---------------------------------------------------------------
	// Joins：用关联表字段做过滤
	// SQL: SELECT orders.*, User.* FROM orders
	//      LEFT JOIN users User ON orders.user_id = User.id
	//      WHERE User.name = 'bob';
	// ---------------------------------------------------------------
	fmt.Println("\n===== Joins：按用户名查订单 =====")
	bobOrders, err := orderDAO.ListByUserName(ctx, "bob")
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		for _, o := range bobOrders {
			fmt.Printf("  订单号: %d, 商品: %s, 用户: %s\n",
				o.OrderNo, o.ProductName, o.User.Name)
		}
	}

	// ---------------------------------------------------------------
	// 组合：Scopes + Preload（分页查已支付订单）
	// ---------------------------------------------------------------
	fmt.Println("\n===== 组合：分页查已支付订单 =====")
	paidOrders, total, err := orderDAO.ListPaidWithUser(ctx, 1, 10)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		fmt.Printf("已支付订单共 %d 条:\n", total)
		for _, o := range paidOrders {
			fmt.Printf("  %s → %s, %d分, 订单号: %d\n",
				o.User.Name, o.ProductName, o.Amount, o.OrderNo)
		}
	}

	// ---------------------------------------------------------------
	// Preload vs Joins 选择指南
	// ---------------------------------------------------------------
	fmt.Println("\n===== 选择指南 =====")
	fmt.Println("Preload：")
	fmt.Println("  ✓ 一对多关联（用户→订单列表）")
	fmt.Println("  ✓ 不需要用关联字段做 WHERE")
	fmt.Println("  ✓ 关联数据量大时，IN 查询比 JOIN 效率更高")
	fmt.Println("")
	fmt.Println("Joins：")
	fmt.Println("  ✓ 需要用关联字段做 WHERE/ORDER BY")
	fmt.Println("  ✓ 一对一关联")
	fmt.Println("  ✓ 需要单次 SQL 完成（减少网络往返）")

	fmt.Println("\n===== 关联查询演示完成 =====")
}
