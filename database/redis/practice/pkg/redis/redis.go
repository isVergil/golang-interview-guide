package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"redis-practice/internal/config"
)

// NewRedis 根据配置初始化 go-redis 客户端并验证连接
func NewRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	opt := &redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	// 连接池
	if cfg.PoolSize > 0 {
		opt.PoolSize = cfg.PoolSize
	}
	if cfg.MinIdleConns > 0 {
		opt.MinIdleConns = cfg.MinIdleConns
	}
	if cfg.MaxIdleConns > 0 {
		opt.MaxIdleConns = cfg.MaxIdleConns
	}
	if cfg.ConnMaxLifetime > 0 {
		opt.ConnMaxLifetime = time.Duration(cfg.ConnMaxLifetime) * time.Second
	}
	if cfg.ConnMaxIdleTime > 0 {
		opt.ConnMaxIdleTime = time.Duration(cfg.ConnMaxIdleTime) * time.Second
	}

	// 超时
	if cfg.DialTimeout > 0 {
		opt.DialTimeout = time.Duration(cfg.DialTimeout) * time.Second
	}
	if cfg.ReadTimeout > 0 {
		opt.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	}
	if cfg.WriteTimeout > 0 {
		opt.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second
	}

	rdb := redis.NewClient(opt)

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return rdb, nil
}

// Close 关闭 Redis 客户端连接
func Close(rdb *redis.Client) error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
