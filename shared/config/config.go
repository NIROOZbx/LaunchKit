package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	HTTPAddr     string        `mapstructure:"http_addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string

	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MinOpenConns    int           `mapstructure:"min_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxIdleTime     time.Duration `mapstructure:"max_idle_time"`

	DSN string
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}


type LogConfig struct {
	Level      string `mapstructure:"level"`
	File       string `mapstructure:"file"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAgeDays int    `mapstructure:"max_age_days"`
}

type JwtConfig struct {
	AccessExpiryMinutes int    `mapstructure:"access_expiry_minutes"`
	RefreshExpiryHours  int    `mapstructure:"refresh_expiry_hours"`
	AccessTokenSecret   string `mapstructure:"access_token_secret"`
	RefreshTokenSecret  string `mapstructure:"refresh_token_secret"`
	Environment         string `mapstructure:"environment"`
	FrontendURL         string `mapstructure:"frontend_url"`
}

type Config struct {
	Core     ServerConfig   `mapstructure:"core"`
	Gateway  ServerConfig   `mapstructure:"gateway"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Jwt     JwtConfig   `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {

	yamlPath := os.Getenv("CONFIG_PATH")
	if yamlPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH env var is required")
	}

	v := viper.New()
	v.SetConfigFile(yamlPath)
	v.AutomaticEnv()

	v.BindEnv("database.user", "POSTGRES_USER")
	v.BindEnv("database.password", "POSTGRES_PASSWORD")
	v.BindEnv("database.host", "POSTGRES_HOST")
	v.BindEnv("database.port", "POSTGRES_PORT")
	v.BindEnv("database.name", "POSTGRES_DB")

	v.BindEnv("redis.addr", "REDIS_ADDR")

	v.BindEnv("jwt.access_token_secret", "ACCESS_SECRET")
	v.BindEnv("jwt.refresh_token_secret", "REFRESH_SECRET")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.Database.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	log.Println("Loaded all configs ✅")
	return &cfg, nil
}

func validate(cfg *Config) error {
	rules := []struct {
		value string
		name  string
	}{
		{cfg.Database.User,     "POSTGRES_USER"},
		{cfg.Database.Password, "POSTGRES_PASSWORD"},
		{cfg.Database.Host,     "POSTGRES_HOST"},
		{cfg.Database.Name,     "POSTGRES_DB"},
		{cfg.Database.Port,     "POSTGRES_PORT"},
		{cfg.Redis.Addr,        "REDIS_ADDR"},
		{cfg.Jwt.AccessTokenSecret, "ACCESS_SECRET"},
		{cfg.Jwt.RefreshTokenSecret, "REFRESH_SECRET"},
	}

	for _, rule := range rules {
		if rule.value == "" {
			return fmt.Errorf("required config missing: %s", rule.name)
		}
	}
	return nil
}