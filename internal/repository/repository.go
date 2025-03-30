package repository

import (
	"context"
	"log/slog"

	"github.com/denisEMPS/gRPC_file_service/internal/domain"
)

type Image interface {
	SaveImage(ctx context.Context, imageData []byte, imageName string) error
	GetImage(ctx context.Context, imageName string) ([]byte, error)
	ListImages(ctx context.Context) ([]domain.ImageInfo, error)
}

type Repository struct {
	Image
}

func New(storageDir string, log *slog.Logger) *Repository {
	return &Repository{Image: NewImageFileSystem(storageDir, log)}
}
