package app

import (
	"github.com/Cheasezz/fileService/config"
	grpcsrv "github.com/Cheasezz/fileService/internal/grpc"
	"github.com/Cheasezz/fileService/internal/repo"
	"github.com/Cheasezz/fileService/internal/service"
	"github.com/Cheasezz/fileService/pkg/logger"
)

type App struct {
	GRPC *grpcsrv.App
	l    logger.Logger
}

func New(l logger.Logger, cfg *config.Config) *App {
	const op = "app.New"
	log := l.With("op", op)

	db, err := repo.New()
	if err != nil {
		log.Error("Cant init file system db: %v", err)
		panic("app create error")
	}

	service := service.New(db, l)

	grpcApp := grpcsrv.New(l, cfg.GRPC, service)

	return &App{GRPC: grpcApp, l: l}
}

func (a *App) Close() {
	const op = "app.Close"
	log := a.l.With("op", op)

	log.Info("stopping gRPC server")
	a.GRPC.Close()
}
