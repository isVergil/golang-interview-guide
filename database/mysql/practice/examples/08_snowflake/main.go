// 08_snowflake：雪花算法分布式 ID 演示
//
// 64 位 ID = 1 bit 符号 + 41 bits 时间戳 + 10 bits 机器ID + 12 bits 序列号
// 特点：趋势递增、全局唯一、不依赖 DB、高性能（单节点 400w+/s）
//
// 运行：cd practice && go run examples/08_snowflake/main.go
package main

import (
	"context"
	"fmt"
	"sync"

	"mysql-practice/examples/setup"
	"mysql-practice/internal/dao"
	"mysql-practice/internal/model"
	"mysql-practice/pkg/database"
	"mysql-practice/pkg/snowflake"
)

func main() {
	db := setup.MustSetup()
	defer database.Close(db)

	orderDAO := dao.NewOrderDAO(db)
	ctx := context.Background()

	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引和外键冲突
	// 先删 orders（有外键引用 users），再删 users
	db.Unscoped().Where("1 = 1").Delete(&model.Order{})
	db.Unscoped().Where("1 = 1").Delete(&model.User{})

	// 创建测试用户
	user := &model.User{Name: "test_user", Age: 25, Email: "snowflake_test@test.com"}
	db.Create(user)

	// ---------------------------------------------------------------
	// 一、基本使用
	// ---------------------------------------------------------------
	fmt.Println("===== 雪花 ID 基本使用 =====")

	node, _ := snowflake.NewNode(1) // 机器 ID = 1

	for i := 0; i < 5; i++ {
		id := node.Generate()
		ts, nodeID, seq := snowflake.ParseID(id)
		fmt.Printf("ID: %d\n  时间: %s\n  节点: %d\n  序列: %d\n\n",
			id, ts.Format("2006-01-02 15:04:05.000"), nodeID, seq)
	}

	// ---------------------------------------------------------------
	// 二、双 ID 模式：自增主键 + 雪花业务号
	//
	// 主键(ID):   自增 BIGINT，保证 B+Tree 顺序写入，内部使用
	// 订单号(OrderNo): 雪花 ID，全局唯一，对外暴露
	// ---------------------------------------------------------------
	fmt.Println("===== 双 ID 模式：订单使用雪花业务号 =====")

	orderNo := node.Generate()
	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      user.ID,
		ProductName: "雪花算法测试商品",
		Amount:      9900,
		Status:      0,
	}
	orderDAO.Create(ctx, order)

	// 用业务订单号查询（对外接口用 OrderNo）
	saved, _ := orderDAO.GetByOrderNo(ctx, orderNo)
	fmt.Printf("自增主键 ID:  %d  (内部用，不对外暴露)\n", saved.ID)
	fmt.Printf("雪花订单号:   %d  (对外用，全局唯一)\n", saved.OrderNo)
	fmt.Printf("商品: %s, 用户: %s\n", saved.ProductName, saved.User.Name)

	ts, nodeID, seq := snowflake.ParseID(saved.OrderNo)
	fmt.Printf("订单号反解: 时间=%s, 节点=%d, 序列=%d\n",
		ts.Format("15:04:05.000"), nodeID, seq)

	// ---------------------------------------------------------------
	// 三、并发安全测试
	// ---------------------------------------------------------------
	fmt.Println("\n===== 并发安全测试 =====")

	const total = 10000
	ids := make([]int64, total)
	var wg sync.WaitGroup

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ids[idx] = node.Generate()
		}(i)
	}
	wg.Wait()

	// 检查唯一性
	seen := make(map[int64]bool, total)
	duplicates := 0
	for _, id := range ids {
		if seen[id] {
			duplicates++
		}
		seen[id] = true
	}
	fmt.Printf("生成 %d 个 ID，重复: %d\n", total, duplicates)

	// 检查趋势递增
	ordered := 0
	for i := 1; i < len(ids); i++ {
		if ids[i] > ids[i-1] {
			ordered++
		}
	}
	fmt.Printf("趋势递增比例: %.1f%%\n", float64(ordered)/float64(total-1)*100)

	// ---------------------------------------------------------------
	// 四、多节点
	// ---------------------------------------------------------------
	fmt.Println("\n===== 多节点生成 =====")

	node1, _ := snowflake.NewNode(1)
	node2, _ := snowflake.NewNode(2)
	node3, _ := snowflake.NewNode(3)

	id1 := node1.Generate()
	id2 := node2.Generate()
	id3 := node3.Generate()

	fmt.Printf("节点1: %d\n", id1)
	fmt.Printf("节点2: %d\n", id2)
	fmt.Printf("节点3: %d\n", id3)
	fmt.Println("→ 不同节点生成的 ID 全局唯一且不冲突")

	// ---------------------------------------------------------------
	// 五、分布式 ID 方案对比
	// ---------------------------------------------------------------
	fmt.Println("\n===== 分布式 ID 方案对比 =====")
	fmt.Println("方案             优点                     缺点")
	fmt.Println("─────────────────────────────────────────────────────────────")
	fmt.Println("DB 自增          简单                     单点瓶颈、分库后不唯一")
	fmt.Println("UUID             无依赖                   无序、索引性能差、占空间")
	fmt.Println("雪花算法         有序、高性能、无DB依赖   依赖时钟、需分配机器ID")
	fmt.Println("Leaf(美团)       高可用                   需部署额外服务")
	fmt.Println("Redis INCR       简单、有序               依赖Redis、有持久化风险")

	fmt.Println("\n===== 雪花算法演示完成 =====")
}
