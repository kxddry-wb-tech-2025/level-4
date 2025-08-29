package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"stresser/internal/delivery/client"
	"stresser/internal/examples"
	"stresser/internal/logging"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logging.New("dev")
	defer logger.Sync()

	victim := os.Getenv("VICTIM")
	if victim == "" {
		panic("VICTIM environment variable is not set")
	}

	stresser := client.NewStresser(victim, ctx)
	fmt.Println("Stresser started")

	stresser.Stress(examples.Orders, true, runtime.GOMAXPROCS(0))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ch
		cancel()
	}()

	logger.Listen(stresser.Logs())

	<-ctx.Done()
	time.Sleep(time.Second)
}
