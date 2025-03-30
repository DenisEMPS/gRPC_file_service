package app

import (
	"log/slog"

	grpcapp "github.com/denisEMPS/gRPC_file_service/internal/app/grpc"
	"github.com/denisEMPS/gRPC_file_service/internal/repository"
	"github.com/denisEMPS/gRPC_file_service/internal/service"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storageDir string) *App {
	repo := repository.New(storageDir, log)
	service := service.New(repo, log)

	grpcApp := grpcapp.New(log, grpcPort, service)

	return &App{
		GRPCServer: grpcApp,
	}
}
