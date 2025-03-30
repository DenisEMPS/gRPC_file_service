package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/denisEMPS/gRPC_file_service/internal/domain"
	"github.com/denisEMPS/gRPC_file_service/internal/repository"
)

var (
	ErrImageAlreadyExists = errors.New("image already exists")
	ErrImageIsNotExists   = errors.New("image does not exists")
)

type ImageService struct {
	repo repository.Image
	log  *slog.Logger
}

func NewImageService(repo repository.Image, log *slog.Logger) *ImageService {
	return &ImageService{repo: repo, log: log}
}

func (s *ImageService) UploadImage(ctx context.Context, imageData []byte, imageName string) error {
	const op = "image_service.UploadImage"

	log := s.log.With(
		slog.String("op", op),
		slog.String("image", imageName),
	)

	err := s.repo.SaveImage(ctx, imageData, imageName)
	if err != nil {
		if errors.Is(err, repository.ErrImageAlreadyExists) {
			log.Warn("failed to upload image", slog.String("error", err.Error()))
			return ErrImageAlreadyExists
		}
		log.Error("failed to upload image", slog.String("error", err.Error()))
		return err
	}

	log.Info("Image uploaded successfully", slog.String("image", imageName))

	return nil
}

func (s *ImageService) DownloadImage(ctx context.Context, imageName string) ([]byte, error) {
	const op = "image_service.DownloadImage"

	log := s.log.With(
		slog.String("op", op),
		slog.String("image", imageName),
	)

	imageData, err := s.repo.GetImage(ctx, imageName)
	if err != nil {
		if errors.Is(err, repository.ErrImageIsNotExists) {
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

	images, err := s.repo.ListImages(ctx)
	if err != nil {
		log.Error("failed to list images", slog.String("error", err.Error()))
		return nil, err
	}

	return images, nil
}
