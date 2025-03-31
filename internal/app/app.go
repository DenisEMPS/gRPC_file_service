package app

import (
	"log/slog"

	grpcapp "github.com/denisEMPS/gRPC_file_service/internal/app/grpc"
	"github.com/denisEMPS/gRPC_file_service/internal/service/image"
	"github.com/denisEMPS/gRPC_file_service/internal/storage/disk"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storageDir string) *App {
	storage := disk.NewImage(storageDir, log)
	service := image.New(storage, log)
	grpcApp := grpcapp.New(log, grpcPort, service)

	return &App{
		GRPCServer: grpcApp,
	}
}
