// 06_optimistic_lock：乐观锁并发控制
//
// 乐观锁 vs 悲观锁：
//
//	乐观锁：不加锁，用 version 字段检测冲突，适合冲突率低的场景
//	悲观锁：SELECT FOR UPDATE 加行锁，阻塞并发，适合冲突率高的场景
//
// 运行：cd practice && go run examples/06_optimistic_lock/main.go
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

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

	// 准备商品：库存 10
	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引（SKU）导致冲突
	db.Unscoped().Where("1 = 1").Delete(&model.Product{})
	product := &model.Product{SKU: "PHONE-001", Name: "iPhone 16", Price: 799900, Stock: 10}
	db.Create(product)
	fmt.Printf("初始库存: %d, version: %d\n", product.Stock, product.Version)

	// ---------------------------------------------------------------
	// 演示1：正常扣减
	// ---------------------------------------------------------------
	fmt.Println("\n===== 正常乐观锁扣减 =====")

	p, _ := productDAO.GetBySKU(ctx, "PHONE-001")
	err := productDAO.DeductStockOptimistic(ctx, p.ID, 2, p.Version)
	if err != nil {
		fmt.Printf("扣减失败: %v\n", err)
	} else {
		p, _ = productDAO.GetBySKU(ctx, "PHONE-001")
		fmt.Printf("扣减 2 件成功: stock=%d, version=%d\n", p.Stock, p.Version)
	}

	// ---------------------------------------------------------------
	// 演示2：模拟冲突（用旧 version 去更新）
	// ---------------------------------------------------------------
	fmt.Println("\n===== 模拟版本冲突 =====")

	// 用户A 读取 version=1
	oldVersion := p.Version

	// 用户B 先扣减了一次，version 变成 2
	productDAO.DeductStockOptimistic(ctx, p.ID, 1, p.Version)
	p, _ = productDAO.GetBySKU(ctx, "PHONE-001")
	fmt.Printf("用户B 扣减成功: stock=%d, version=%d\n", p.Stock, p.Version)

	// 用户A 用旧 version=1 去扣减，失败
	err = productDAO.DeductStockOptimistic(ctx, p.ID, 1, oldVersion)
	if err != nil {
		fmt.Printf("用户A 扣减失败(预期): %v\n", err)
	}

	// ---------------------------------------------------------------
	// 演示3：带重试的乐观锁（生产环境标准写法）
	// ---------------------------------------------------------------
	fmt.Println("\n===== 带重试的乐观锁 =====")

	err = deductWithRetry(ctx, productDAO, p.ID, 1, 3)
	if err != nil {
		fmt.Printf("重试后仍失败: %v\n", err)
	} else {
		p, _ = productDAO.GetBySKU(ctx, "PHONE-001")
		fmt.Printf("重试扣减成功: stock=%d, version=%d\n", p.Stock, p.Version)
	}

	// ---------------------------------------------------------------
	// 演示4：并发扣减（10 个 goroutine 各扣 1 件）
	// ---------------------------------------------------------------
	fmt.Println("\n===== 并发扣减测试 =====")

	// 重置库存为 10
	db.Model(&model.Product{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"stock": 10, "version": 0,
	})

	var wg sync.WaitGroup
	var successCount int64
	var failCount int64

	for i := 0; i < 20; i++ { // 20 个并发，但只有 10 个库存
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			err := deductWithRetry(ctx, productDAO, p.ID, 1, 5)
			if err != nil {
				atomic.AddInt64(&failCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}
	wg.Wait()

	p, _ = productDAO.GetBySKU(ctx, "PHONE-001")
	fmt.Printf("20 并发扣减结果: 成功=%d, 失败=%d, 剩余库存=%d\n",
		successCount, failCount, p.Stock)
	fmt.Println("→ 乐观锁保证库存不会变成负数")

	fmt.Println("\n===== 乐观锁演示完成 =====")
}

// deductWithRetry 带重试的乐观锁扣减（生产环境标准模式）
func deductWithRetry(ctx context.Context, productDAO *dao.ProductDAO, id uint, quantity int, maxRetry int) error {
	for i := 0; i < maxRetry; i++ {
		// 每次重试都重新读取最新 version
		p, err := productDAO.GetBySKU(ctx, "PHONE-001")
		if err != nil {
			return err
		}
		err = productDAO.DeductStockOptimistic(ctx, id, quantity, p.Version)
		if err == nil {
			return nil // 成功
		}
		// 冲突，继续重试
	}
	return fmt.Errorf("optimistic lock failed after %d retries", maxRetry)
}
