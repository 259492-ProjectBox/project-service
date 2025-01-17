package repositories

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type FileExtensionRepository interface {
	repository[models.FileExtension]
	GetFileExtension(ctx context.Context, tx *gorm.DB, file *multipart.FileHeader) (*models.FileExtension, error)
}

type fileExtensionRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.FileExtension]
}

func NewFileExtensionRepository(db *gorm.DB) FileExtensionRepository {
	return &fileExtensionRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.FileExtension](db),
	}
}

func (r *fileExtensionRepositoryImpl) GetFileExtension(ctx context.Context, tx *gorm.DB, file *multipart.FileHeader) (*models.FileExtension, error) {
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "unknown"
	}
	fileExtension := &models.FileExtension{}
	if err := tx.WithContext(ctx).Where("mime_type = ?", fileType).First(fileExtension).Error; err != nil {
		return nil, fmt.Errorf("file type not found: %w", err)
	}

	return fileExtension, nil
}
