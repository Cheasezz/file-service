package grpcsrv

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/Cheasezz/fileService/internal/service"
	"github.com/Cheasezz/fileService/pkg/logger"
	file "github.com/Cheasezz/fileService/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout"`
}

type server struct {
	file.UnimplementedFileServer
	service *service.Service
}

type App struct {
	server *grpc.Server
	log    logger.Logger
	port   int
}

func New(l logger.Logger, cfg Config, s *service.Service) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			l.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainStreamInterceptor(
		recovery.StreamServerInterceptor(recoveryOpts...),
		logging.StreamServerInterceptor(InterceptorLogger(l), loggingOpts...),
	))

	file.RegisterFileServer(gRPCServer, &server{service: s})

	return &App{
		log:    l,
		server: gRPCServer,
		port:   cfg.Port,
	}
}

func InterceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, int(lvl), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcsrv.Run"
	log := a.log.With("op", op)

	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", "addr", l.Addr().String())

	if err = a.server.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Close() {
	a.server.GracefulStop()
}
