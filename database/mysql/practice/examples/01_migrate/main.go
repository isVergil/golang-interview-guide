package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"

	"mysql-practice/examples/setup"
	"mysql-practice/pkg/database"
)

// 01_migrate：建库建表最佳实践
//
// 生产环境不用 AutoMigrate，用版本化 SQL 迁移
// 本文件演示企业级 DDL 规范和 GORM 迁移技巧
func main() {
	db := setup.MustSetup()
	defer database.Close(db)

	ctx := context.Background()

	// ---------------------------------------------------------------
	// 一、查看 AutoMigrate 生成的建表 DDL
	// 验证 GORM 生成的表结构是否符合预期
	// ---------------------------------------------------------------
	fmt.Println("===== 查看建表 DDL =====")

	tables := []string{"users", "orders", "products"}
	for _, table := range tables {
		var createSQL string
		row := db.WithContext(ctx).Raw("SHOW CREATE TABLE " + table).Row()
		var tableName string
		if err := row.Scan(&tableName, &createSQL); err != nil {
			log.Printf("show create table %s: %v", table, err)
			continue
		}
		fmt.Printf("\n-- %s\n%s;\n", table, createSQL)
	}

	// ---------------------------------------------------------------
	// 二、企业级建表规范速查
	//
	// 1. 字符集：utf8mb4（支持 emoji），排序规则 utf8mb4_general_ci
	// 2. 存储引擎：InnoDB（事务 + 行锁 + MVCC）
	// 3. 主键：自增 BIGINT 或雪花 ID（禁止用 UUID 做聚簇索引，B+Tree 分裂严重）
	// 4. 字段注释：每个字段必须有 COMMENT（GORM tag: comment:xxx）
	// 5. NOT NULL：字段尽量 NOT NULL + DEFAULT，避免 NULL 导致索引失效
	// 6. 金额：用 INT/BIGINT 存分，不用 DECIMAL/FLOAT
	// 7. 时间：DATETIME 或 TIMESTAMP，不要用字符串
	// 8. 索引：
	//    - 区分度高的字段放前面（复合索引遵循最左前缀）
	//    - 避免在索引列上做函数运算（WHERE YEAR(created_at) = 2024 走不了索引）
	//    - 用 EXPLAIN 验证
	// 9. 软删除：deleted_at 字段 + 索引
	// ---------------------------------------------------------------

	// ---------------------------------------------------------------
	// 三、手动创建索引（AutoMigrate 之外的补充索引）
	// 生产环境通常在 SQL migration 文件里写
	// ---------------------------------------------------------------
	fmt.Println("\n\n===== 手动创建复合索引 =====")

	// 复合索引示例：按年龄+余额联合查询
	// MySQL 没有 CREATE INDEX IF NOT EXISTS，需要先查是否存在
	var indexCount int64
	db.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'users' AND index_name = 'idx_age_balance'",
	).Scan(&indexCount)

	var result *gorm.DB
	if indexCount == 0 {
		result = db.WithContext(ctx).Exec("CREATE INDEX idx_age_balance ON users(age, balance)")
	} else {
		fmt.Println("索引 idx_age_balance 已存在，跳过")
	}
	if result.Error != nil {
		log.Printf("create index: %v", result.Error)
	} else {
		fmt.Println("创建复合索引 idx_age_balance(age, balance) 成功")
	}

	// ---------------------------------------------------------------
	// 四、EXPLAIN 分析查询是否走索引
	// ---------------------------------------------------------------
	fmt.Println("\n===== EXPLAIN 分析 =====")

	type ExplainResult struct {
		ID           int    `gorm:"column:id"`
		SelectType   string `gorm:"column:select_type"`
		Table        string `gorm:"column:table"`
		Type         string `gorm:"column:type"`
		PossibleKeys string `gorm:"column:possible_keys"`
		Key          string `gorm:"column:key"`
		Rows         int    `gorm:"column:rows"`
		Extra        string `gorm:"column:Extra"`
	}

	queries := []struct {
		desc string
		sql  string
	}{
		{"主键查询", "EXPLAIN SELECT * FROM users WHERE id = 1"},
		{"唯一索引查询", "EXPLAIN SELECT * FROM users WHERE email = 'test@example.com'"},
		{"复合索引查询", "EXPLAIN SELECT * FROM users WHERE age = 25 AND balance > 1000"},
		{"索引失效（函数）", "EXPLAIN SELECT * FROM users WHERE UPPER(email) = 'TEST@EXAMPLE.COM'"},
	}

	for _, q := range queries {
		fmt.Printf("\n-- %s\n", q.desc)
		var results []ExplainResult
		db.WithContext(ctx).Raw(q.sql).Scan(&results)
		for _, r := range results {
			fmt.Printf("  type=%-6s key=%-20s rows=%-5d extra=%s\n",
				r.Type, r.Key, r.Rows, r.Extra)
		}
	}

	// ---------------------------------------------------------------
	// 五、生产环境推荐的迁移工具
	//
	// 1. golang-migrate/migrate  (最流行，SQL 文件驱动)
	//    migrate -path ./migrations -database "mysql://..." up
	//
	// 2. pressly/goose (支持 Go 函数 + SQL 双模式)
	//    goose -dir ./migrations mysql "user:pass@/dbname" up
	//
	// 迁移文件命名规范：
	//   000001_create_users.up.sql
	//   000001_create_users.down.sql
	//   000002_add_user_phone.up.sql
	//   000002_add_user_phone.down.sql
	// ---------------------------------------------------------------
	fmt.Println("\n\n===== 建表最佳实践演示完成 =====")
	fmt.Println("提示：生产环境请用 golang-migrate 或 goose 管理 DDL 变更")
}
