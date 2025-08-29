package main

import (
	"context"
	"log"
	"metrics/internal/delivery/http"
	"metrics/internal/metrics"
	"time"
)

const MetricsDuration = 2 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	server := http.New()
	go metrics.RecordMetrics(ctx, MetricsDuration)

	if err := server.Start(":8080"); err != nil {
		log.Printf("Failed to start server: %v", err)
	}

	cancel()
	time.Sleep(MetricsDuration)
}
