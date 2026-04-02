package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`

	// 连接池
	MaxConns        int `yaml:"max_conns"`
	MinConns        int `yaml:"min_conns"`
	MaxConnLifetime int `yaml:"max_conn_lifetime"`
	MaxConnIdleTime int `yaml:"max_conn_idle_time"`

	// 超时（秒）
	ConnectTimeout int `yaml:"connect_timeout"`
}

// DSN 返回 PostgreSQL 连接字符串
func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&connect_timeout=%d",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode, c.ConnectTimeout)
}

// LoadConfig 从 yaml 文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	content := resolveEnvVars(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := cfg.Postgres.validate(); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}

	return &cfg, nil
}

// resolveEnvVars 解析 ${ENV_VAR:default} 格式的环境变量
func resolveEnvVars(content string) string {
	result := content
	for {
		start := strings.Index(result, "${")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}")
		if end == -1 {
			break
		}
		end += start

		expr := result[start+2 : end]
		parts := strings.SplitN(expr, ":", 2)
		envKey := parts[0]
		defaultVal := ""
		if len(parts) == 2 {
			defaultVal = parts[1]
		}

		envVal := os.Getenv(envKey)
		if envVal == "" {
			envVal = defaultVal
		}

		result = result[:start] + envVal + result[end+1:]
	}
	return result
}

func (c *PostgresConfig) validate() error {
	if c.Host == "" {
		return fmt.Errorf("postgres.host is required")
	}
	if c.Port == 0 {
		return fmt.Errorf("postgres.port is required")
	}
	if c.User == "" {
		return fmt.Errorf("postgres.user is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("postgres.dbname is required")
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 5
	}
	return nil
}
