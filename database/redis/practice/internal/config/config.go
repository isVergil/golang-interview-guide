package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis RedisConfig `yaml:"redis"`
}

type RedisConfig struct {
	// 连接信息
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`

	// 连接池
	PoolSize        int `yaml:"pool_size"`
	MinIdleConns    int `yaml:"min_idle_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int `yaml:"conn_max_idle_time"`

	// 超时（秒）
	DialTimeout  int `yaml:"dial_timeout"`
	ReadTimeout  int `yaml:"read_timeout"`
	WriteTimeout int `yaml:"write_timeout"`
}

// Addr 返回 host:port 格式地址
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LoadConfig 从 yaml 文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	// 支持环境变量覆盖：${ENV_VAR:default}
	content := resolveEnvVars(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := cfg.Redis.validate(); err != nil {
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

func (c *RedisConfig) validate() error {
	if c.Host == "" {
		return fmt.Errorf("redis.host is required")
	}
	if c.Port == 0 {
		return fmt.Errorf("redis.port is required")
	}
	return nil
}
