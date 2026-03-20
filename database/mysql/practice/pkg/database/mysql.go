package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mysql-practice/internal/config"
)

// NewMySQL 根据配置初始化 GORM 连接
func NewMySQL(cfg *config.MySQLConfig) (*gorm.DB, error) {
	// GORM 日志级别
	logLevel := parseLogLevel(cfg.LogLevel)

	// 慢查询阈值
	slowThreshold := time.Duration(cfg.SlowThreshold) * time.Millisecond
	if slowThreshold == 0 {
		slowThreshold = 200 * time.Millisecond
	}

	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// 建立连接
	db, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	// 获取底层 sql.DB 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
	}

	// 验证连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return db, nil
}

// Close 关闭数据库连接
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// parseLogLevel 将配置字符串转为 GORM 日志级别
func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
