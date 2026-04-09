package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	App      App      `mapstructure:"app"`
	HTTP     HTTP     `mapstructure:"http"`
	GRPC     GRPC     `mapstructure:"grpc"`
	Postgres Postgres `mapstructure:"postgres"`
	Kafka    Kafka    `mapstructure:"kafka"`
	Redis    Redis    `mapstructure:"redis"`
	OTLP     OTLP     `mapstructure:"otlp"`
	Telegram Telegram `mapstructure:"telegram"`
}

type App struct {
	Name     string `mapstructure:"name"`
	Env      Env    `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
}

type HTTP struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (h HTTP) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type GRPC struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (g GRPC) Addr() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func (p Postgres) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.DB, p.SSLMode,
	)
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

type Redis struct {
	Addr string `mapstructure:"addr"`
}

type OTLP struct {
	Endpoint string `mapstructure:"endpoint"`
}

type Telegram struct {
	Token string `mapstructure:"token"`
}

// Load reads configs in order: common.yaml → {env}.yaml → sensitive.yaml → env vars.
// Each layer overrides only the fields it specifies.
func Load(configDir string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// 1. common.yaml — shared defaults
	v.SetConfigName("common")
	v.AddConfigPath(configDir)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read common.yaml: %w", err)
	}

	// 2. {env}.yaml — environment-specific overrides
	env := v.GetString("app.env")
	if env == "" {
		env = string(EnvLocal)
	}
	v.SetConfigName(env)
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("read %s.yaml: %w", env, err)
	}

	// 3. sensitive.yaml — local secrets, gitignored
	v.SetConfigName("sensitive")
	_ = v.MergeInConfig() // optional, ignore if missing

	// 4. env vars — highest priority
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
