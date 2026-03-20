package setup

import (
	"log"
	"os"

	"gorm.io/gorm"

	"mysql-practice/internal/config"
	"mysql-practice/internal/model"
	"mysql-practice/pkg/database"
)

// MustSetup 初始化数据库连接 + 自动建表，失败直接退出
// 所有 example 共用，避免每个文件重复写初始化代码
//
// 使用方式：
//
//	db := setup.MustSetup()
//	defer database.Close(db)
func MustSetup() *gorm.DB {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("[setup] load config: %v", err)
	}

	db, err := database.NewMySQL(&cfg.MySQL)
	if err != nil {
		log.Fatalf("[setup] init mysql: %v", err)
	}

	// 自动建表（仅开发环境使用，生产用 migration）
	if err := db.AutoMigrate(&model.User{}, &model.Order{}, &model.Product{}); err != nil {
		log.Fatalf("[setup] auto migrate: %v", err)
	}

	return db
}
