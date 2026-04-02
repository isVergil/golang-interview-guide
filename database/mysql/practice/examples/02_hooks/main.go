package main

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"mysql-practice/examples/setup"
	"mysql-practice/internal/dao"
	"mysql-practice/internal/model"
	"mysql-practice/pkg/database"
)

// 02_hooks：GORM Hooks 生命周期演示
//
// Hooks 是 GORM 在 Create/Update/Delete/Query 前后自动调用的方法
// 定义在 Model 上，适合做数据清洗、自动填充、审计日志
func main() {
	db := setup.MustSetup()
	defer database.Close(db)

	userDAO := dao.NewUserDAO(db)
	ctx := context.Background()

	// ---------------------------------------------------------------
	// BeforeCreate Hook 已定义在 model/user.go 中：
	//   - Email 自动转小写 + 去空格
	//   - Name 自动去空格
	// ---------------------------------------------------------------

	fmt.Println("===== Hooks：BeforeCreate 自动清洗 =====")

	// 测试数据清理用物理删除（Unscoped），避免软删除记录占唯一索引导致后续插入冲突
	cleanupUser(db, "alice@example.com")
	cleanupUser(db, "bob@test.com")

	// 故意传入带空格和大写的数据
	user := &model.User{Name: "  Alice  ", Age: 25, Email: "  Alice@Example.COM  "}
	if err := userDAO.Create(ctx, user); err != nil {
		fmt.Printf("创建失败: %v\n", err)
		return
	}

	fmt.Printf("传入  Name=[ Alice ], Email=[  Alice@Example.COM  ]\n")
	fmt.Printf("存储  Name=[%s], Email=[%s]\n", user.Name, user.Email)
	fmt.Println("→ BeforeCreate 自动 TrimSpace + ToLower")

	// ---------------------------------------------------------------
	// BeforeUpdate Hook 同样生效
	// ---------------------------------------------------------------
	fmt.Println("\n===== Hooks：BeforeUpdate 自动清洗 =====")

	// 用 struct 更新时触发 Hook（注意：map 更新不触发 Hook）
	user.Name = "  BOB  "
	user.Email = "  BOB@TEST.COM  "
	db.Save(user)

	// 重新查出来验证
	updated, _ := userDAO.GetByID(ctx, user.ID)
	fmt.Printf("更新  Name=[ BOB ], Email=[  BOB@TEST.COM  ]\n")
	fmt.Printf("存储  Name=[%s], Email=[%s]\n", updated.Name, updated.Email)
	fmt.Println("→ BeforeUpdate 自动清洗")

	// ---------------------------------------------------------------
	// 注意事项
	// ---------------------------------------------------------------
	fmt.Println("\n===== 重要：Hook 不触发的情况 =====")
	fmt.Println("1. db.Updates(map[string]interface{}{...})  → 不触发 Hook")
	fmt.Println("2. db.Exec(\"UPDATE ...\")                   → 不触发 Hook")
	fmt.Println("3. db.UpdateColumn / UpdateColumns           → 不触发 Hook")
	fmt.Println("4. 批量操作 db.Where(...).Updates(...)       → 不触发 Hook")
	fmt.Println("")
	fmt.Println("只有通过 db.Save(&model) 或 db.Create(&model) 才会触发")

	// 清理
	db.Unscoped().Where("id = ?", user.ID).Delete(&model.User{})
	fmt.Println("\n===== Hooks 演示完成 =====")
}

func cleanupUser(db *gorm.DB, email string) {
	// Unscoped 物理删除，确保唯一索引释放
	db.Unscoped().Where("email = ?", strings.ToLower(strings.TrimSpace(email))).Delete(&model.User{})
}
