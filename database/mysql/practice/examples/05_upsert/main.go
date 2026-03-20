// 05_upsert：幂等写入（INSERT ON DUPLICATE KEY UPDATE）
//
// 场景：同步外部数据、消息消费去重、配置批量导入
// 按唯一键冲突时更新，不重复插入
//
// 运行：cd practice && go run examples/05_upsert/main.go
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

	productDAO := dao.NewProductDAO(db)
	ctx := context.Background()

	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引（SKU）导致冲突
	db.Unscoped().Where("1 = 1").Delete(&model.Product{})

	// ---------------------------------------------------------------
	// 单条 Upsert
	// ---------------------------------------------------------------
	fmt.Println("===== 单条 Upsert =====")

	// 第一次插入
	p1 := &model.Product{SKU: "BOOK-001", Name: "Go 语言圣经", Price: 6800, Stock: 100}
	if err := productDAO.Upsert(ctx, p1); err != nil {
		fmt.Printf("upsert failed: %v\n", err)
		return
	}
	product, _ := productDAO.GetBySKU(ctx, "BOOK-001")
	fmt.Printf("第1次: name=%s, price=%d, stock=%d\n", product.Name, product.Price, product.Stock)

	// 第二次：相同 SKU，更新价格和库存
	// SQL: INSERT INTO products (...) VALUES (...)
	//      ON DUPLICATE KEY UPDATE name=VALUES(name), price=VALUES(price), stock=VALUES(stock)
	p2 := &model.Product{SKU: "BOOK-001", Name: "Go 语言圣经(第2版)", Price: 8800, Stock: 200}
	if err := productDAO.Upsert(ctx, p2); err != nil {
		fmt.Printf("upsert failed: %v\n", err)
		return
	}
	product, _ = productDAO.GetBySKU(ctx, "BOOK-001")
	fmt.Printf("第2次: name=%s, price=%d, stock=%d\n", product.Name, product.Price, product.Stock)
	fmt.Println("→ SKU 冲突，自动更新而非报错")

	// ---------------------------------------------------------------
	// 批量 Upsert
	// ---------------------------------------------------------------
	fmt.Println("\n===== 批量 Upsert =====")

	products := []model.Product{
		{SKU: "BOOK-001", Name: "Go 语言圣经(第3版)", Price: 9900, Stock: 300},
		{SKU: "BOOK-002", Name: "Redis 设计与实现", Price: 5500, Stock: 50},
		{SKU: "BOOK-003", Name: "MySQL 技术内幕", Price: 7900, Stock: 80},
	}
	if err := productDAO.BatchUpsert(ctx, products, 100); err != nil {
		fmt.Printf("batch upsert failed: %v\n", err)
		return
	}

	// 验证
	for _, sku := range []string{"BOOK-001", "BOOK-002", "BOOK-003"} {
		p, _ := productDAO.GetBySKU(ctx, sku)
		fmt.Printf("  %s: name=%s, price=%d, stock=%d\n", p.SKU, p.Name, p.Price, p.Stock)
	}
	fmt.Println("→ BOOK-001 更新为第3版，BOOK-002/003 新插入")

	// ---------------------------------------------------------------
	// 使用场景
	// ---------------------------------------------------------------
	fmt.Println("\n===== 典型场景 =====")
	fmt.Println("1. 外部数据同步：定时拉取第三方商品数据，按 SKU 幂等写入")
	fmt.Println("2. 消息消费去重：消费 MQ 消息时，按业务唯一键幂等处理")
	fmt.Println("3. 配置批量导入：Excel 导入配置，存在则更新，不存在则新增")
	fmt.Println("4. 数据修复脚本：批量修复数据，不怕重复执行")

	fmt.Println("\n===== Upsert 演示完成 =====")
}
