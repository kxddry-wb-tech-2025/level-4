package config

import "time"

type Config struct {
	Env  string      `yaml:"env" env-default:"dev"`
	SMTP EmailConfig `yaml:"smtp"`
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
