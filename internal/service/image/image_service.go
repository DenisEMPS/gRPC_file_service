package image

import (
	"context"
	"errors"
	"log/slog"

	"github.com/denisEMPS/gRPC_file_service/internal/domain"
	"github.com/denisEMPS/gRPC_file_service/internal/storage/disk"
)

type ImageStorage interface {
	Save(ctx context.Context, imageData []byte, imageName string) error
	Get(ctx context.Context, imageName string) ([]byte, error)
	List(ctx context.Context) ([]domain.ImageInfo, error)
}

var (
	ErrImageAlreadyExists = errors.New("image already exists")
	ErrImageIsNotExists   = errors.New("image does not exists")
)

type ImageService struct {
	storage ImageStorage
	log     *slog.Logger
}

func New(storage ImageStorage, log *slog.Logger) *ImageService {
	return &ImageService{storage: storage, log: log}
}

func (s *ImageService) UploadImage(ctx context.Context, imageData []byte, imageName string) error {
	const op = "image_service.UploadImage"

	log := s.log.With(
		slog.String("op", op),
		slog.String("image", imageName),
	)

	err := s.storage.Save(ctx, imageData, imageName)
	if err != nil {
		if errors.Is(err, disk.ErrImageAlreadyExists) {
			log.Warn("failed to upload image", slog.String("error", err.Error()))
			return ErrImageAlreadyExists
		}
		log.Error("failed to upload image", slog.String("error", err.Error()))
		return err
	}

	log.Info("image uploaded successfully", slog.String("image", imageName))

	return nil
}

func (s *ImageService) DownloadImage(ctx context.Context, imageName string) ([]byte, error) {
	const op = "image_service.DownloadImage"

	log := s.log.With(
		slog.String("op", op),
		slog.String("image", imageName),
	)

	imageData, err := s.storage.Get(ctx, imageName)
	if err != nil {
		if errors.Is(err, disk.ErrImageIsNotExists) {
			log.Warn("failed to download image", slog.String("error", err.Error()))
			return nil, ErrImageIsNotExists
		}
		log.Error("failed to download image", slog.String("error", err.Error()))
		return nil, err
	}

	return imageData, nil
}

func (s *ImageService) ListImage(ctx context.Context) ([]domain.ImageInfo, error) {
	const op = "image_service.ListImage"

	log := s.log.With(
		slog.String("op", op),
	)

	images, err := s.storage.List(ctx)
	if err != nil {
		log.Error("failed to list images", slog.String("error", err.Error()))
		return nil, err
	}

	return images, nil
}
