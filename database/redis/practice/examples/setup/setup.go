package setup

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"

	"redis-practice/internal/config"
	pkgRedis "redis-practice/pkg/redis"
)

// MustSetup 初始化 Redis 连接，失败直接退出
// 所有 example 共用，避免每个文件重复写初始化代码
//
// 使用方式：
//
//	rdb := setup.MustSetup()
//	defer pkgRedis.Close(rdb)
func MustSetup() *redis.Client {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("[setup] load config: %v", err)
	}

	rdb, err := pkgRedis.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("[setup] init redis: %v", err)
	}

	return rdb
}
