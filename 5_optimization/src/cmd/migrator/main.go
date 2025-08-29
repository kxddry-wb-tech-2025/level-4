package main

import (
	"database/sql"
	"errors"
	"fmt"
	"l0/internal/config"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"

	// migration
	"github.com/golang-migrate/migrate/v4"

	// drivers
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrationConfig struct {
	St            config.Storage `yaml:"storage" env-required:"true"` // use dbname = postgres here
	Operation     string         `env:"OPERATION" yaml:"operation" env-default:"up"`
	DbsMigrations []entry        `yaml:"dbs_migrations" env-required:"true"`
}

type entry struct {
	Name string `yaml:"name"` // DBName
	Path string `yaml:"path"` // Path for migrations
}

func main() {
	confPath := os.Getenv("CONFIG_PATH")
	if confPath == "" {
		panic("CONFIG_PATH env variable not set")
	}

	var cfg migrationConfig
	if err := cleanenv.ReadConfig(confPath, &cfg); err != nil {
		panic(err)
	}
	op := cfg.Operation

	if len(cfg.DbsMigrations) == 0 {
		panic("No dbs migrations configured")
	}

	for _, m := range cfg.DbsMigrations {
		name, path := m.Name, m.Path
		ccfg := cfg.St
		ccfg.DBName = name
		err := ensureDBexists(name, ccfg)
		if err != nil {
			panic(err)
		}
		link := link(ccfg)

		doOneMigration(link, path, op)
	}

	fmt.Println("migration successful")
}

func doOneMigration(link, path, op string) {
	m, err := migrate.New("file://"+path, link)
	if err != nil {
		panic(err)
	}
	switch {
	case op == "" || op == "up":
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Nothing to migrate at", path)
				return
			}
			panic(err)
		}
		return
	case op == "down":
		if err = m.Force(1); err != nil {
			panic(err)
		}
		if err = m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Nothing to migrate at", path)
				return
			}
			panic(err)
		}
	default:
		panic("Unknown operation: " + op)
	}
}

func link(cfg config.Storage) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
}

func dataSourceName(cfg config.Storage) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func ensureDBexists(dbname string, adminCfg config.Storage) error {
	adminCfg.DBName = "postgres"
	_db, err := sql.Open("postgres", dataSourceName(adminCfg))
	if err != nil {
		return err
	}
	defer _db.Close()

	_, err = _db.Exec("CREATE DATABASE " + dbname)
	if err != nil && !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "no change") {
		return err
	}
	return nil
}
