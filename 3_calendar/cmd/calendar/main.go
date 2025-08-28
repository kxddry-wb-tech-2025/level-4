package main

import (
	"calendar/internal/config"
	"calendar/internal/delivery/http"
	notifyDelivery "calendar/internal/delivery/notify"
	"calendar/internal/delivery/smtp"
	"calendar/internal/lib/fanin"
	"calendar/internal/logging"
	eventsSvc "calendar/internal/service/events"
	notifySvc "calendar/internal/service/notify"
	workerSvc "calendar/internal/service/worker"
	redisStore "calendar/internal/storage/redis"
	"calendar/internal/storage/txmanager"
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
	defer logger.Sync()

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redis := redisStore.NewStorage(cfg.Redis)

	emailClient := smtp.NewEmailClient(cfg.SMTP)
	sender := notifyDelivery.NewSender(emailClient, cfg.SMTP)

	txmgr, err := txmanager.New(mainCtx, cfg.Storage.DSN())
	if err != nil {
		panic(err)
	}

	w := workerSvc.NewWorker(redis, txmgr.AsWorkerTxManager(), sender)
	wlogs := w.Logs()

	go w.Handle(mainCtx)

	jobs := make(chan any, 100)

	eSvc := eventsSvc.NewService(mainCtx, txmgr.AsEventsTxManager(), jobs)
	elogs := eSvc.Logs()
	nSvc := notifySvc.NewService(mainCtx, w)
	nlogs := nSvc.Logs()

	go nSvc.Process(jobs)

	srv := http.NewServer(mainCtx, cfg.Server, eSvc)
	logs := srv.Logs()

	// Fan in all logs into the logger
	logger.Listen(fanin.FanIn(mainCtx, logs, elogs, nlogs, wlogs))

	if err := srv.Start(); err != nil {
	}
}
