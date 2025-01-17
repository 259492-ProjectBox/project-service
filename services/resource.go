package services

import (
	"context"
	"mime/multipart"

	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ResourceService interface {
	UploadAssetResource(ctx context.Context, file *multipart.FileHeader, description string, programId int) (*models.AssetResource, error)
	GetAssetResourceByProgramID(ctx context.Context, programId string) ([]models.AssetResource, error)
	DeleteAssetResourceByID(ctx context.Context, id string) error
	GetDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error)
	DeleteProjectResourceByID(ctx context.Context, id string, filePath *string) error
	GetResourcesByProjectID(ctx context.Context, projectID string) ([]models.Resource, error)
}

type resourceService struct {
	resourceRepository repositories.ResourceRepository
}

func NewResourceService(resourceRepository repositories.ResourceRepository) ResourceService {
	return &resourceService{
		resourceRepository: resourceRepository,
	}
}

func (s *resourceService) UploadAssetResource(ctx context.Context, file *multipart.FileHeader, description string, programId int) (*models.AssetResource, error) {
	assetResource := &models.AssetResource{
		Description: description,
		ProgramID:   programId,
	}
	return s.resourceRepository.CreateAssetResource(ctx, file, assetResource)
}

func (s *resourceService) GetDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error) {
	return s.resourceRepository.FindDetailedResourceByID(ctx, id)
}

func (s *resourceService) DeleteProjectResourceByID(ctx context.Context, id string, filePath *string) error {
	return s.resourceRepository.DeleteProjectResourceByID(ctx, id, filePath)
}

func (s *resourceService) GetAssetResourceByProgramID(ctx context.Context, programId string) ([]models.AssetResource, error) {
	return s.resourceRepository.FindAssetResourcesByProgramID(ctx, programId)
}
func (s *resourceService) DeleteAssetResourceByID(ctx context.Context, id string) error {
	return s.resourceRepository.DeleteAssetResourceByID(ctx, id)
}

func (s *resourceService) GetResourcesByProjectID(ctx context.Context, projectID string) ([]models.Resource, error) {
	return s.resourceRepository.FindByProjectID(ctx, projectID)
}
