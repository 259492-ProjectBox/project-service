package repositories

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
)

type UploadRepository interface {
	UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader) error
	DeleteUploadedFiles(ctx context.Context, objectNames []string)
}

type uploadRepositoryImpl struct {
	minioClient *minio.Client
}

func NewUploadRepository(minioClient *minio.Client) UploadRepository {
	return &uploadRepositoryImpl{
		minioClient: minioClient,
	}
}

func (r *uploadRepositoryImpl) UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = r.minioClient.PutObject(ctx, os.Getenv("MINIO_PROJECT_BUCKET"), objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *uploadRepositoryImpl) DeleteUploadedFiles(ctx context.Context, objectNames []string) {
	for _, objectName := range objectNames {
		err := r.minioClient.RemoveObject(ctx, os.Getenv("MINIO_PROJECT_BUCKET"), objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Failed to delete file from MinIO: %s, error: %v", objectName, err)
		}
	}
}
