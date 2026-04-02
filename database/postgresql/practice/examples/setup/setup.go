package setup

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"pg-practice/internal/config"
	pkgPG "pg-practice/pkg/postgres"
)

// MustSetup 初始化 PostgreSQL 连接池，失败直接退出
// 所有 example 共用，避免每个文件重复写初始化代码
func MustSetup() *pgxpool.Pool {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("[setup] load config: %v", err)
	}

	pool, err := pkgPG.NewPostgres(&cfg.Postgres)
	if err != nil {
		log.Fatalf("[setup] init postgres: %v", err)
	}

	return pool
}
