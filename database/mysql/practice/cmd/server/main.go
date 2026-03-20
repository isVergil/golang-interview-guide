package main

import (
	"fmt"
	"log"

	"mysql-practice/internal/config"
	"mysql-practice/internal/model"
	"mysql-practice/pkg/database"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	fmt.Println("[1] 配置加载成功")

	// 2. 初始化数据库
	db, err := database.NewMySQL(&cfg.MySQL)
	if err != nil {
		log.Fatalf("init mysql: %v", err)
	}
	defer database.Close(db)
	fmt.Println("[2] 数据库连接成功")

	// 3. 自动建表
	if err := db.AutoMigrate(&model.User{}, &model.Order{}, &model.Product{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}
	fmt.Println("[3] 自动建表完成（users, orders, products）")

	fmt.Println("\n基建就绪，可运行 examples/ 下的各个示例")
}
