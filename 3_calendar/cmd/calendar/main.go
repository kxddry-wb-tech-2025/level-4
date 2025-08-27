package main

import (
	"calendar/internal/config"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config.yaml"
	}

	cfg := config.MustLoad(path)
	_ = cfg
}
