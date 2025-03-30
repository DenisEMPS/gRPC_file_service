package service

import (
	"context"
	"log/slog"

	"github.com/denisEMPS/gRPC_file_service/internal/domain"
	"github.com/denisEMPS/gRPC_file_service/internal/repository"
)

type Image interface {
	UploadImage(ctx context.Context, imageData []byte, imageName string) error
	DownloadImage(ctx context.Context, imageName string) ([]byte, error)
	ListImage(ctx context.Context) ([]domain.ImageInfo, error)
}

type Service struct {
	Image
}

func New(repo *repository.Repository, log *slog.Logger) *Service {
	return &Service{Image: NewImageService(repo.Image, log)}
}
