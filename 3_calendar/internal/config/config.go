package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string        `yaml:"env" env-default:"dev"`
	Server  ServerConfig  `yaml:"server"`
	SMTP    EmailConfig   `yaml:"smtp"`
	Storage StorageConfig `yaml:"storage"`
	Redis   RedisConfig   `yaml:"redis"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	Database int    `yaml:"database" env-default:"0"`
}

type EmailConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `env:"SMTP_PASSWORD" env-required:"true"`
}

type ServerConfig struct {
	Port        int           `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type StorageConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `env:"STORAGE_PASSWORD" env-required:"true"`
	Database string `yaml:"database" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

func (sc *StorageConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		sc.Host, sc.Port, sc.Username, sc.Password, sc.Database, sc.SSLMode)
}

func MustLoad(path string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
