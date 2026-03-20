package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
}

type MySQLConfig struct {
	// 连接信息
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
	Charset   string `yaml:"charset"`
	Timezone  string `yaml:"timezone"`
	ParseTime bool   `yaml:"parse_time"`

	// 连接池
	MaxOpenConns    int `yaml:"max_open_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int `yaml:"conn_max_idle_time"`

	// 日志
	LogLevel      string `yaml:"log_level"`
	SlowThreshold int    `yaml:"slow_threshold"`
}

// DSN 拼接 GORM 连接所需的 Data Source Name
func (c *MySQLConfig) DSN() string {
	charset := c.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	timezone := c.Timezone
	if timezone == "" {
		timezone = "Local"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName,
		charset, c.ParseTime, timezone)
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

	if err := cfg.MySQL.validate(); err != nil {
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

		expr := result[start+2 : end] // ENV_VAR:default
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

// validate 校验必填字段
func (c *MySQLConfig) validate() error {
	if c.Host == "" {
		return fmt.Errorf("mysql.host is required")
	}
	if c.Port == 0 {
		return fmt.Errorf("mysql.port is required")
	}
	if c.User == "" {
		return fmt.Errorf("mysql.user is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("mysql.dbname is required")
	}
	return nil
}
