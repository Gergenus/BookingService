package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MinioImageRepository struct {
	minioClient *minio.Client
	bucketName  string
	endpoint    string
}

type ImageRepositoryInterface interface {
	AddImage(ctx context.Context, image *multipart.FileHeader) (string, error)
	DeleteImage(ctx context.Context, objectName string) error
}

func NewMinioImageRepository(minioClient *minio.Client, bucketName string, endpoint string) *MinioImageRepository {
	return &MinioImageRepository{
		minioClient: minioClient,
		bucketName:  bucketName,
		endpoint:    endpoint,
	}
}

// returns url
func (m *MinioImageRepository) AddImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	const op = "image_repository.AddImage"
	fileToAdd, err := image.Open()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	_, err = m.minioClient.PutObject(ctx, m.bucketName, objectName, fileToAdd, image.Size, minio.PutObjectOptions{
		ContentType: ".jpg",
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return objectName, nil
}

func (m *MinioImageRepository) DeleteImage(ctx context.Context, objectName string) error {
	const op = "image_repository.DeleteImage"
	err := m.minioClient.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
