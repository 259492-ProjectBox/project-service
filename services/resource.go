package services

// import (
// 	"context"
// 	"errors"

// 	"github.com/project-box/models"
// 	"github.com/project-box/repositories"
// )

// type ResourceService interface {
// 	CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error)
// 	GetResourceByID(ctx context.Context, id string) (*models.Resource, error)
// 	DeleteResource(ctx context.Context, id string) error
// 	GetResourcesByProjectID(ctx context.Context, projectID string) ([]models.Resource, error)
// }

// type resourceService struct {
// 	resourceRepository repositories.ResourceRepository
// }

// func NewResourceService(resourceRepository repositories.ResourceRepository) ResourceService {
// 	return &resourceService{
// 		resourceRepository: resourceRepository,
// 	}
// }

// func (s *resourceService) CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error) {
// 	if resource.ProjectID == 0 {
// 		return nil, errors.New("project_id is required")
// 	}
// 	if resource.ResourceTypeID == 0 {
// 		return nil, errors.New("resource_type_id is required")
// 	}
// 	if resource.URL == "" {
// 		return nil, errors.New("url is required")
// 	}

// 	return s.resourceRepository.CreateResource(ctx, resource)
// }

// func (s *resourceService) GetResourceByID(ctx context.Context, id string) (*models.Resource, error) {
// 	return s.resourceRepository.FindResourceByID(ctx, id)
// }

// func (s *resourceService) DeleteResource(ctx context.Context, id string) error {
// 	return s.resourceRepository.DeleteResource(ctx, id)
// }

// func (s *resourceService) GetResourcesByProjectID(ctx context.Context, projectID string) ([]models.Resource, error) {
// 	return s.resourceRepository.FindByProjectID(ctx, projectID)
// }
