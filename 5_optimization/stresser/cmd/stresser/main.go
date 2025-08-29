package main

import (
	"context"
	"os"
	"os/signal"
	"stresser/internal/delivery/client"
	"stresser/internal/examples"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	victim := os.Getenv("VICTIM")
	if victim == "" {
		logger.Fatal("VICTIM environment variable is not set")
	}

	logger.Info("Starting stresser", zap.String("victim", victim))

	stresser := client.NewStresser(victim, ctx)
	stresser.Stress(examples.Orders, true)

	logger.Info("Stress test completed, processing logs")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ch
		cancel()
	}()

	for entry := range stresser.Logs() {
		logger.Info("log entry", zap.Any("entry", entry.Message))
	}

	logger.Info("Stresser finished")
}
