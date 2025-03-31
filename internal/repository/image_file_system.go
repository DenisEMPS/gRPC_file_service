package imagefs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/denisEMPS/gRPC_file_service/internal/domain"
	"golang.org/x/sys/unix"
)

var (
	ErrImageAlreadyExists = errors.New("image already exists")
	ErrImageIsNotExists   = errors.New("images does not exists")
)

type ImageFS struct {
	storageDir string
	logger     *slog.Logger
}

func New(storageDir string, logger *slog.Logger) *ImageFS {
	if err := os.Mkdir(storageDir, 0755); err != nil {
		log.Fatal("failed to create image directory: ", err.Error())
	}

	return &ImageFS{storageDir: storageDir, logger: logger}
}

func (r *ImageFS) SaveImage(ctx context.Context, imageData []byte, imageName string) error {
	const op = "image_file_system.SaveImage"

	imagePath := filepath.Join(r.storageDir, imageName)

	if _, err := os.Stat(imagePath); !os.IsNotExist(err) {
		return fmt.Errorf("%s: %w", op, ErrImageAlreadyExists)
	}

	err := os.WriteFile(imagePath, imageData, 0644)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ImageFS) GetImage(ctx context.Context, imageName string) ([]byte, error) {
	const op = "image_file_system.GetImage"

	imagePath := filepath.Join(r.storageDir, imageName)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %w", op, ErrImageIsNotExists)
	}

	imageData, err := os.ReadFile(imagePath)
	if err != nil || len(imageData) == 0 {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return imageData, nil
}

func (r *ImageFS) ListImages(ctx context.Context) ([]domain.ImageInfo, error) {
	const op = "image_file_system.ListImages"

	var imagesInfo []domain.ImageInfo

	err := filepath.Walk(r.storageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || info.Name()[0] == '.' {
			return nil
		}

		createTime, err := r.GetFileCreationTime(path)
		if err != nil {
			return err
		}
		imagesInfo = append(imagesInfo, domain.ImageInfo{
			Name:      info.Name(),
			CreatedAt: createTime,
			UpdatedAt: info.ModTime(),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return imagesInfo, nil
}

func (r *ImageFS) GetFileCreationTime(filePath string) (time.Time, error) {
	var statx unix.Statx_t
	err := unix.Statx(unix.AT_FDCWD, filePath, 0, unix.STATX_BTIME, &statx)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get create time for file %s: %w", filePath, err)
	}

	return time.Unix(statx.Btime.Sec, int64(statx.Btime.Nsec)), nil
}
