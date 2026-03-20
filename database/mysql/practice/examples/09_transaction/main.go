// 09_transaction：事务与转账（悲观锁 FOR UPDATE）
//
// 演示 GORM 事务的标准用法：
//   - db.Transaction() 闭包：返回 nil 自动 Commit，返回 error 自动 Rollback
//   - FOR UPDATE 悲观锁：防止并发读到脏余额
//   - gorm.Expr 原子更新：避免读-改-写竞态
//   - DAO.WithTx(tx)：保证事务内操作走同一个连接
//
// 运行：cd practice && go run examples/09_transaction/main.go
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"mysql-practice/examples/setup"
	"mysql-practice/internal/dao"
	"mysql-practice/internal/model"
	"mysql-practice/internal/service"
	"mysql-practice/pkg/database"
)

func main() {
	db := setup.MustSetup()
	defer database.Close(db)

	userDAO := dao.NewUserDAO(db)
	userService := service.NewUserService(db, userDAO)
	ctx := context.Background()

	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引导致后续插入冲突
	db.Unscoped().Where("1 = 1").Delete(&model.User{})

	// ---------------------------------------------------------------
	// 一、准备转账双方
	// ---------------------------------------------------------------
	fmt.Println("===== 准备用户 =====")

	alice, _ := userService.CreateUser(ctx, "alice", 25, "alice_tx@test.com")
	bob, _ := userService.CreateUser(ctx, "bob", 30, "bob_tx@test.com")

	// 充值
	userDAO.Update(ctx, alice.ID, map[string]interface{}{"balance": 10000}) // 100 元
	userDAO.Update(ctx, bob.ID, map[string]interface{}{"balance": 5000})    // 50 元

	alice, _ = userDAO.GetByID(ctx, alice.ID)
	bob, _ = userDAO.GetByID(ctx, bob.ID)
	fmt.Printf("alice: balance=%d (100元)\n", alice.Balance)
	fmt.Printf("bob:   balance=%d (50元)\n", bob.Balance)

	// ---------------------------------------------------------------
	// 二、正常转账
	// ---------------------------------------------------------------
	fmt.Println("\n===== 正常转账：alice → bob 30元 =====")

	err := userService.Transfer(ctx, alice.ID, bob.ID, 3000)
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		alice, _ = userDAO.GetByID(ctx, alice.ID)
		bob, _ = userDAO.GetByID(ctx, bob.ID)
		fmt.Printf("转账成功\n")
		fmt.Printf("alice: %d → %d\n", 10000, alice.Balance)
		fmt.Printf("bob:   %d → %d\n", 5000, bob.Balance)
	}

	// ---------------------------------------------------------------
	// 三、转账失败 → 自动回滚
	// ---------------------------------------------------------------
	fmt.Println("\n===== 余额不足：alice → bob 999元 =====")

	beforeAlice, _ := userDAO.GetByID(ctx, alice.ID)
	beforeBob, _ := userDAO.GetByID(ctx, bob.ID)

	err = userService.Transfer(ctx, alice.ID, bob.ID, 99900)
	if err != nil {
		fmt.Printf("预期失败: %v\n", err)
	}

	afterAlice, _ := userDAO.GetByID(ctx, alice.ID)
	afterBob, _ := userDAO.GetByID(ctx, bob.ID)
	fmt.Printf("alice 余额: %d → %d (不变，已回滚)\n", beforeAlice.Balance, afterAlice.Balance)
	fmt.Printf("bob   余额: %d → %d (不变，已回滚)\n", beforeBob.Balance, afterBob.Balance)

	// ---------------------------------------------------------------
	// 四、并发转账测试（FOR UPDATE 保证一致性）
	// ---------------------------------------------------------------
	fmt.Println("\n===== 并发转账测试 =====")

	// 重置余额
	userDAO.Update(ctx, alice.ID, map[string]interface{}{"balance": 10000})
	userDAO.Update(ctx, bob.ID, map[string]interface{}{"balance": 10000})
	fmt.Println("重置: alice=10000, bob=10000")

	// 10 个 goroutine 同时 alice → bob 转 10 元（1000分）
	var wg sync.WaitGroup
	var successCount int64
	var failCount int64

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := userService.Transfer(ctx, alice.ID, bob.ID, 1000)
			if err != nil {
				atomic.AddInt64(&failCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}()
	}
	wg.Wait()

	alice, _ = userDAO.GetByID(ctx, alice.ID)
	bob, _ = userDAO.GetByID(ctx, bob.ID)
	fmt.Printf("10次并发转账(每次10元): 成功=%d, 失败=%d\n", successCount, failCount)
	fmt.Printf("alice: %d, bob: %d\n", alice.Balance, bob.Balance)
	fmt.Printf("总和: %d (应为 20000，资金守恒)\n", alice.Balance+bob.Balance)

	// ---------------------------------------------------------------
	// 五、事务要点总结
	// ---------------------------------------------------------------
	fmt.Println("\n===== 事务要点 =====")
	fmt.Println("1. db.Transaction(func(tx *gorm.DB) error { ... })")
	fmt.Println("   返回 nil → Commit，返回 error → Rollback")
	fmt.Println("")
	fmt.Println("2. FOR UPDATE 悲观锁")
	fmt.Println("   tx.Set(\"gorm:query_option\", \"FOR UPDATE\").First(&user, id)")
	fmt.Println("   锁住行，防止并发读到旧余额")
	fmt.Println("")
	fmt.Println("3. gorm.Expr 原子更新")
	fmt.Println("   Updates(map{\"balance\": gorm.Expr(\"balance - ?\", amount)})")
	fmt.Println("   SQL 层面 SET balance = balance - 3000，不依赖 Go 层的值")
	fmt.Println("")
	fmt.Println("4. DAO.WithTx(tx)")
	fmt.Println("   事务内创建临时 DAO，保证所有操作走同一个 DB 连接")

	fmt.Println("\n===== 事务演示完成 =====")
}
