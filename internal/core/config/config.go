package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB     DBConfig
	HTTP   HTTPConfig
	Logger LoggerConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-default:"marketplace_db"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

type HTTPConfig struct {
	Port    string        `env:"HTTP_PORT" env-default:"8080"`
	Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"10s"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

func MustLoadConfig() *Config {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = ".env"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, reading from env", configPath)
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			panic("failed to read env: " + err.Error())
		}
	} else {
		log.Printf("Loading config from %s", configPath)
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			panic("failed to read config: " + err.Error())
		}
	}

	return &cfg
}
