package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config — базовая конфигурация сервиса (.env или config.yml, без дефолтов).
type Config struct {
	Port         string `mapstructure:"port"`
	GrpcPort     string `mapstructure:"grpc_port"`
	Env          string `mapstructure:"env"`
	LogLevel     string `mapstructure:"log_level"`
	LogFormat    string `mapstructure:"log_format"`
	LogAddSource bool   `mapstructure:"log_add_source"`
}

// ErrNoConfigFile — конфигурационный файл (.env или config.yml) не найден.
var ErrNoConfigFile = errors.New("config: no config file found (.env or config.yml)")

// Load читает конфигурацию: сначала .env, при отсутствии — config.yml (или CONFIG_FILE). Без дефолтов; при отсутствии файла или обязательных полей возвращает ошибку.
func Load() (Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	_ = v.BindEnv("port", "PORT")
	_ = v.BindEnv("grpc_port", "GRPC_PORT")
	_ = v.BindEnv("env", "ENV")
	_ = v.BindEnv("log_level", "LOG_LEVEL")
	_ = v.BindEnv("log_format", "LOG_FORMAT")
	_ = v.BindEnv("log_add_source", "LOG_ADD_SOURCE")

	var read bool
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, fmt.Errorf("config: read .env: %w", err)
		}
	} else {
		read = true
		// .env часто использует UPPERCASE; маппим в lowercase для mapstructure
		for _, pair := range [][]string{{"PORT", "port"}, {"GRPC_PORT", "grpc_port"}, {"ENV", "env"}, {"LOG_LEVEL", "log_level"}, {"LOG_FORMAT", "log_format"}, {"LOG_ADD_SOURCE", "log_add_source"}} {
			if v.IsSet(pair[0]) {
				v.Set(pair[1], v.Get(pair[0]))
			}
		}
	}
	if !read {
		configFile := os.Getenv("CONFIG_FILE")
		if configFile == "" {
			configFile = "config.yml"
		}
		v.SetConfigFile(configFile)
		v.SetConfigType("yaml")
		if err := v.ReadInConfig(); err != nil {
			var notFound viper.ConfigFileNotFoundError
			if errors.As(err, &notFound) {
				return Config{}, ErrNoConfigFile
			}
			return Config{}, fmt.Errorf("config: read %s: %w", configFile, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("config: unmarshal: %w", err)
	}
	if cfg.Port == "" {
		return Config{}, errors.New("config: port is required")
	}
	if cfg.GrpcPort == "" {
		return Config{}, errors.New("config: grpc_port is required")
	}
	if cfg.Env == "" {
		return Config{}, errors.New("config: env is required")
	}
	if cfg.LogLevel == "" {
		return Config{}, errors.New("config: log_level is required")
	}
	if cfg.LogFormat == "" {
		return Config{}, errors.New("config: log_format is required")
	}
	return cfg, nil
}
