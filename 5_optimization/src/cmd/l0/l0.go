package main

import (
	"context"
	"errors"
	"l0/internal/config"
	"l0/internal/handlers"
	"l0/internal/storage/cache"
	"l0/internal/storage/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	initCfg "github.com/kxddry/go-utils/pkg/config"
	initLog "github.com/kxddry/go-utils/pkg/logger"
	"github.com/kxddry/go-utils/pkg/logger/handlers/sl"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var cfg config.Config
	initCfg.MustParseConfig(&cfg)
	log := initLog.SetupLogger(cfg.Env)
	log.Debug("debug enabled")
	st, err := postgres.NewStorage(cfg.Storage)
	if err != nil {
		panic(err)
	}
	cacher := cache.NewCache(cfg.Cache.TTL, cfg.Cache.Limit)

	err = handlers.LoadCache(ctx, cacher, st)
	if err != nil {
		panic(err)
	}

	// I'm too lazy to refactor this whole thing. Hence I'm not creating a Server layer.
	// I'm just going to remove Kafka here and accept POST requests with order data.
	e := echo.New()
	e.Validator = NewValidator()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet}, // only GET allowed
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: cfg.Server.Timeout}))

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      e,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	e.GET("/order/:id", handlers.GetOrderHandler(st, cacher))
	e.POST("/order", handlers.CreateOrderHandler(st))

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to start", sl.Err(err))
	}

	// graceful shutdown
	down := make(chan os.Signal, 1)
	signal.Notify(down, os.Signal(syscall.SIGTERM), os.Signal(syscall.SIGINT))
	<-down
	log.Info("shutting down")

	_ = srv.Shutdown(ctx)

}

type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
}
