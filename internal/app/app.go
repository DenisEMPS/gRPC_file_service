package app

import (
	"log/slog"

	grpcapp "github.com/denisEMPS/gRPC_file_service/internal/app/grpc"
	imagefs "github.com/denisEMPS/gRPC_file_service/internal/repository"
	image "github.com/denisEMPS/gRPC_file_service/internal/service"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storageDir string) *App {
	imagefs := imagefs.New(storageDir, log)
	imageservice := image.New(imagefs, log)

	grpcApp := grpcapp.New(log, grpcPort, imageservice)

	return &App{
		GRPCServer: grpcApp,
	}
}
