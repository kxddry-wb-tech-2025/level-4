package config

import "time"

// Config is a structure with configs
type Config struct {
	Env     string  `yaml:"env" env-default:"dev"` // local, dev, prod
	Storage Storage `yaml:"storage"`
	Kafka   Kafka   `yaml:"kafka"`
	Server  Server  `yaml:"server"`
	Cache   Cache   `yaml:"cache"`
}

// Storage is a structure with configs for PostgreSQL
type Storage struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-default:"require"`
}

// Server is a structure with configs for an HTTP server
type Server struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"3s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// Cache is a structure with configs for creating cache
type Cache struct {
	TTL   time.Duration `yaml:"ttl" env-default:"15m"`
	Limit int           `yaml:"limit" env-default:"1000"`
}

// Kafka is a structure with configs for a broker like Kafka
type Kafka struct {
	Brokers []string     `yaml:"brokers" env-required:"true"`
	Reader  ReaderConfig `yaml:"reader" env-required:"true"`
	Writer  WriterConfig `yaml:"writer" env-required:"true"`
}

// ReaderConfig is a structure with config for kafka reader
type ReaderConfig struct {
	Topic          string        `yaml:"topic" env-required:"true"`
	GroupID        string        `yaml:"group_id" env-required:"true"`
	MinBytes       int           `yaml:"min_bytes" env-default:"1"`         // min fetch bytes
	MaxBytes       int           `yaml:"max_bytes" env-default:"1048576"`   // 1MB
	CommitInterval time.Duration `yaml:"commit_interval" env-default:"1s"`  // time.Duration, e.g. 1s
	StartOffset    string        `yaml:"start_offset" env-default:"latest"` // earliest | latest
}

// WriterConfig is a structure with config for kafka writer
type WriterConfig struct {
	Topic           string        `yaml:"topic" env-required:"true"`
	ClientID        string        `yaml:"client_id" env-required:"true"`
	Retries         int           `yaml:"retries" env-default:"5"`
	MaxMessageBytes int           `yaml:"max_message_bytes" env-default:"1048576"`
	Acks            string        `yaml:"acks" env-default:"all"`        // 0 | 1 | all
	Compression     string        `yaml:"compression" env-default:"lz4"` // lz4 | snappy | none | gzip | zstd
	Timeout         time.Duration `yaml:"timeout" env-default:"5s"`      // time.Duration
}
