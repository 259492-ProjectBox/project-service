package services

import (
	"context"
	"github.com/project-box/models"
	"github.com/project-box/repositories"
)

type ResourceService interface {
	GetDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error)
	DeleteProjectResourceByID(ctx context.Context, id string, filePath *string) error
}

type resourceService struct {
	resourceRepository repositories.ResourceRepository
}

func NewResourceService(resourceRepository repositories.ResourceRepository) ResourceService {
	return &resourceService{
		resourceRepository: resourceRepository,
	}
}

func (s *resourceService) GetDetailedResourceByID(ctx context.Context, id string) (*models.DetailedResource, error) {
	return s.resourceRepository.FindDetailedResourceByID(ctx, id)
}

func (s *resourceService) DeleteProjectResourceByID(ctx context.Context, id string, filePath *string) error {
	return s.resourceRepository.DeleteProjectResourceByID(ctx, id, filePath)
}
