package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Cheasezz/fileService/config"
	"github.com/Cheasezz/fileService/internal/app"
	"github.com/Cheasezz/fileService/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("starting application")

	application := app.New(log, cfg)
	defer application.Close()

	go application.GRPC.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Gracefully stopped")
}
