package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the main configuration for the application
type Config struct {
	Env     string         `yaml:"env" env-default:"dev"`
	Server  *ServerConfig  `yaml:"server"`
	SMTP    *EmailConfig   `yaml:"smtp"`
	Storage *StorageConfig `yaml:"storage"`
	Redis   *RedisConfig   `yaml:"redis"`
	Worker  *WorkerConfig  `yaml:"worker"`
}

// RedisConfig is the configuration for the Redis database
type RedisConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	Database int    `yaml:"database" env-default:"0"`
}

// EmailConfig is the configuration for the email server
type EmailConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `env:"SMTP_PASSWORD" env-required:"true"`
}

// ServerConfig is the configuration for the HTTP server
type ServerConfig struct {
	Port        int           `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

// StorageConfig is the configuration for the PostgreSQL database
type StorageConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `env:"STORAGE_PASSWORD" env-required:"true"`
	Database string `yaml:"database" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

// ArchiverConfig is the configuration for the archiver service
type ArchiverConfig struct {
	Interval  time.Duration `yaml:"interval" env-required:"true"`
	OlderThan time.Duration `yaml:"older_than" env-required:"true"`
	BatchSize int           `yaml:"batch_size" env-required:"true"`
}

// WorkerConfig is the configuration for the worker service
type WorkerConfig struct {
	Interval time.Duration `yaml:"interval" env-default:"30s"`
	Limit    int64         `yaml:"limit" env-default:"100"`
}

// DSN returns the Data Source Name for the PostgreSQL database
func (sc *StorageConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		sc.Host, sc.Port, sc.Username, sc.Password, sc.Database, sc.SSLMode)
}

// MustLoad loads the configuration from a file and panics if the file does not exist or the configuration is invalid
func MustLoad(path string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
