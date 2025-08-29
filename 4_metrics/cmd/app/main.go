package main

import (
	"context"
	"log"
	"metrics/internal/delivery/http"
	"metrics/internal/metrics"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const MetricsDuration = 2 * time.Second

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "2112"
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := http.New()
	go metrics.RecordMetrics(ctx, MetricsDuration)

	if err := server.Start(port); err != nil {
		log.Printf("Failed to start server: %v", err)
	}

	cancel()
	time.Sleep(MetricsDuration)
}
