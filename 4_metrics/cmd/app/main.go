package main

import (
	"log"
	"metrics/internal/delivery/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "2112"
	}
	server := http.New()

	if err := server.Start(port); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
