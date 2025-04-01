package filerpc

import (
	"context"
	"errors"

	file_service "github.com/DenisEMPS/proto-repo/gen/go/file-service"

	"github.com/denisEMPS/gRPC_file_service/internal/domain"
	"github.com/denisEMPS/gRPC_file_service/internal/service/image"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ImageService interface {
	UploadImage(ctx context.Context, imageData []byte, imageName string) error
	DownloadImage(ctx context.Context, imageName string) ([]byte, error)
	ListImage(ctx context.Context) ([]domain.ImageInfo, error)
}

type ServerAPI struct {
	file_service.UnimplementedFileServer
	service      ImageService
	semaUpload   chan struct{}
	semaDownload chan struct{}
	semaList     chan struct{}
}

func NewServerApi(imageService ImageService) *ServerAPI {
	return &ServerAPI{
		service:      imageService,
		semaUpload:   make(chan struct{}, 10),
		semaDownload: make(chan struct{}, 10),
		semaList:     make(chan struct{}, 100),
	}
}

func RegisterServer(gRPC *grpc.Server, service ImageService) {
	file_service.RegisterFileServer(gRPC, NewServerApi(service))
}

func (a *ServerAPI) UploadImage(ctx context.Context, req *file_service.UploadImageRequest) (*file_service.UploadImageResponse, error) {
	select {
	case a.semaUpload <- struct{}{}:
		defer func() { <-a.semaUpload }()
	default:
		return nil, status.Error(codes.ResourceExhausted, "too many upload requests, please try again later")
	}

	if len(req.ImageData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request body")
	}

	if req.ImageName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty name field")
	}

	err := a.service.UploadImage(ctx, req.ImageData, req.ImageName)
	if err != nil {
		if errors.Is(err, image.ErrImageAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "image already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &file_service.UploadImageResponse{
		Message: "success",
	}, nil
}

func (a *ServerAPI) DownloadImage(ctx context.Context, req *file_service.DownloadImageRequest) (*file_service.DownloadImageResponse, error) {
	select {
	case a.semaDownload <- struct{}{}:
		defer func() { <-a.semaDownload }()
	default:
		return nil, status.Error(codes.ResourceExhausted, "too many download requests, please try again later")
	}

	if req.ImageName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty name field")
	}

	imageData, err := a.service.DownloadImage(ctx, req.ImageName)
	if err != nil {
		if errors.Is(err, image.ErrImageIsNotExists) {
			return nil, status.Error(codes.NotFound, "image not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &file_service.DownloadImageResponse{
		ImageData: imageData,
	}, nil
}

func (a *ServerAPI) ListImage(ctx context.Context, req *file_service.ListImageRequest) (*file_service.ListImageResponse, error) {
	select {
	case a.semaList <- struct{}{}:
		defer func() { <-a.semaList }()
	default:
		return nil, status.Error(codes.ResourceExhausted, "too many list requests, please try again later")
	}

	imagesInfo, err := a.service.ListImage(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	var grpcImages []*file_service.ListImageResponse_ImageInfo
	for _, image := range imagesInfo {
		createdAt := timestamppb.New(image.CreatedAt)
		updatedAt := timestamppb.New(image.UpdatedAt)

		grpcImages = append(grpcImages, &file_service.ListImageResponse_ImageInfo{
			ImageName: image.Name,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return &file_service.ListImageResponse{
		Images: grpcImages,
	}, nil
}
