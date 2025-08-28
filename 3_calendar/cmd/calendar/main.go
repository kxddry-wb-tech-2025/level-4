package main

import (
	"calendar/internal/config"
	"calendar/internal/delivery/http"
	notifyDelivery "calendar/internal/delivery/notify"
	"calendar/internal/delivery/smtp"
	"calendar/internal/logging"
	eventsSvc "calendar/internal/service/events"
	notifySvc "calendar/internal/service/notify"
	workerSvc "calendar/internal/service/worker"
	redisStore "calendar/internal/storage/redis"
	eventsRepo "calendar/internal/storage/repos/events"
	notifyRepo "calendar/internal/storage/repos/notifications"
	"context"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = os.Setenv("TZ", "UTC")
	_ = godotenv.Load()

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "configs/config.yaml"
	}

	cfg := config.MustLoad(path)
	logger := logging.New(cfg.Env)
	defer logger.Sync() // nolint:errcheck

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Storage
	redis := redisStore.NewStorage(&cfg.Redis)

	// SMTP + sender
	emailClient := smtp.NewEmailClient(&cfg.SMTP)
	sender := notifyDelivery.NewSender(emailClient, &cfg.SMTP)

	// DB Repos
	eRepo, err := eventsRepo.NewRepository(mainCtx, cfg.Storage.DSN())
	if err != nil {
		panic(err)
	}
	nRepo, err := notifyRepo.NewRepository(mainCtx, cfg.Storage.DSN())
	if err != nil {
		panic(err)
	}

	// Worker
	w := workerSvc.NewWorker(redis, sender, nRepo, nil)
	w.Handle(mainCtx)

	// Jobs channel between event service and notify service
	jobs := make(chan any, 100)

	// Services
	eSvc := eventsSvc.NewService(eRepo, jobs)
	nSvc := notifySvc.NewService(nRepo, w)

	// Notify service listens to jobs
	go nSvc.Process(mainCtx, jobs)

	// HTTP server
	srv := http.NewServer(mainCtx, cfg.Server, eSvc)
	logs := srv.Logs()
	logger.Listen(logs)

	if err := srv.Start(); err != nil {
	}
}
