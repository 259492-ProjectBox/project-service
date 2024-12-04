package services

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type UploadService interface {
	UploadObject(ctx context.Context, bucketName string, objectName string, Object io.Reader, fileSize int64, contentType string) (string, error)
	RemoveObject(ctx context.Context, bucketName string, objectName string, opt minio.RemoveObjectOptions) error
}

type uploadServiceImpl struct {
	client *minio.Client
}

func NewMinioService(client *minio.Client) UploadService {
	return &uploadServiceImpl{
		client: client,
	}
}

func (s *uploadServiceImpl) UploadObject(ctx context.Context, bucketName string, objectName string, file io.Reader, fileSize int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, bucketName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	// Generate presigned URL
	url, err := s.client.PresignedGetObject(ctx, bucketName, objectName, time.Hour*24*7, nil)
	if err != nil {
		s.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
		return "", err
	}
	return url.String(), nil
}

func (s *uploadServiceImpl) RemoveObject(ctx context.Context, bucketName string, objectName string, opt minio.RemoveObjectOptions) error {
	return s.client.RemoveObject(context.Background(), bucketName, objectName, opt)
}
