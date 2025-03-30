package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	filerpc "github.com/denisEMPS/gRPC_file_service/internal/grpc/file-rpc"
	"github.com/denisEMPS/gRPC_file_service/internal/service"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	GRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, service *service.Service) *App {
	gRPCServer := grpc.NewServer()

	filerpc.RegisterServer(gRPCServer, service)

	return &App{
		log:        log,
		GRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.GRPCServer.GracefulStop()
}
