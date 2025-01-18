package repositories

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
)

type UploadRepository interface {
	UploadFile(ctx context.Context, bucketName string, objectName string, file *multipart.FileHeader, options minio.PutObjectOptions) error
	DeleteUploadedFiles(ctx context.Context, bucketName string, objectNames []string, removeOptions minio.RemoveObjectOptions) error
	DeleteUploadedFile(ctx context.Context, bucketName string, objectName string, removeOptions minio.RemoveObjectOptions) error
}

type uploadRepositoryImpl struct {
	minioClient *minio.Client
}

func NewUploadRepository(minioClient *minio.Client) UploadRepository {
	return &uploadRepositoryImpl{
		minioClient: minioClient,
	}
}

func (r *uploadRepositoryImpl) UploadFile(ctx context.Context, bucketName string, objectName string, file *multipart.FileHeader, options minio.PutObjectOptions) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = r.minioClient.PutObject(ctx, bucketName, objectName, src, file.Size, options)
	if err != nil {
		return err
	}

	return nil
}

func (r *uploadRepositoryImpl) DeleteUploadedFiles(ctx context.Context, bucketName string, objectNames []string, removeOptions minio.RemoveObjectOptions) error {
	var allErrors []string

	for _, objectName := range objectNames {
		err := r.deleteFile(ctx, bucketName, objectName, removeOptions)
		if err != nil {
			allErrors = append(allErrors, fmt.Sprintf("Failed to delete file %s: %v", objectName, err))
		}
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("errors occurred during deletion: %s", strings.Join(allErrors, "; "))
	}
	return nil
}

func (r *uploadRepositoryImpl) DeleteUploadedFile(ctx context.Context, bucketName string, objectName string, removeOptions minio.RemoveObjectOptions) error {
	return r.deleteFile(ctx, bucketName, objectName, removeOptions)
}

func (r *uploadRepositoryImpl) deleteFile(ctx context.Context, bucketName string, objectName string, removeOptions minio.RemoveObjectOptions) error {
	err := r.minioClient.RemoveObject(ctx, bucketName, objectName, removeOptions)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", objectName, err)
	}
	return nil
}
