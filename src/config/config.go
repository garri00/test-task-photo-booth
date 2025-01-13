package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/viper"
)

type Configs struct {
	Host            string `env:"SERVICE_HOST, default=localhost"`
	Port            string `env:"SERVICE_PORT, default=8080"`
	LogLevel        string `env:"SERVICE_LOGLEVEL, default=DEBUG"`
	Mode            string `env:"SERVICE_MODE, default=TEST"`
	TLS             bool   `env:"SERVICE_TLS, default=false"`
	ViperConfigPath string `env:"SERVICE_CONFIG, required"`
	PostgresConf    PostgresDBConf
	RabbitMQConf    RabbitMQConf
}

// PostgresDBConf creates config for db connection
type PostgresDBConf struct {
	Host     string `env:"SERVICE_PGHOST, required"`
	Port     string `env:"SERVICE_PGPORT, required"`
	Database string `env:"SERVICE_PGDATABASE, required"`
	Username string `env:"SERVICE_PGUSER, required"`
	Password string `env:"SERVICE_PGPASSWORD, required"`
	SSLMode  string `env:"SERVICE_PGSSLMODE, required"`
}

// RabbitMQConf creates config for queue connection
type RabbitMQConf struct {
	Host     string `env:"SERVICE_RMQHOST, required"`
	Port     string `env:"SERVICE_RMQPORT, required"`
	Username string `env:"SERVICE_RMQUSER, required"`
	Password string `env:"SERVICE_RMQPASSWORD, required"`
}

func GetConfig() (Configs, error) {
	var cfg Configs

	if err := godotenv.Load(); err != nil {
		return cfg, fmt.Errorf("loading .env file failed: %w", err)
	}

	ctx := context.Background()
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return cfg, fmt.Errorf("loading env variables failed: %w", err)
	}

	if err := loadViperConfig(cfg.ViperConfigPath); err != nil {
		return cfg, fmt.Errorf("loading viper config failed: %w", err)
	}

	return cfg, nil
}

const (
	viperConfigFileName = "config"
	viperConfigFileType = "json"
)

func loadViperConfig(viperConfigPath string) error {
	viper.SetConfigName(viperConfigFileName) // name of config file (without extension)

	viper.SetConfigType(viperConfigFileType) // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(viperConfigPath)

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("viper.ReadInConfig() failed: %w", err)
	}

	return nil
}
