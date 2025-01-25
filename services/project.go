package services

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	rabbitMQQueue "github.com/project-box/queues/rabbitmq"
	"github.com/project-box/repositories"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type ProjectService interface {
	PublishProjectMessageToElasticSearch(ctx context.Context, action string, projectId int)
	GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error)
	CreateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error)
	UpdateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error)
	DeleteProject(ctx context.Context, id int) error
}

type projectServiceImpl struct {
	rabbitMQChannel *rabbitmq.Channel
	projectRepo     repositories.ProjectRepository
	committeeRepo   repositories.StaffRepository
	programRepo     repositories.ProgramRepository
	courseRepo      repositories.CourseRepository
}

func NewProjectService(
	rabbitMQChannel *rabbitmq.Channel,
	projectRepo repositories.ProjectRepository,
	committeeRepo repositories.StaffRepository,
	programRepo repositories.ProgramRepository,
	courseRepo repositories.CourseRepository,
) ProjectService {
	return &projectServiceImpl{
		rabbitMQChannel: rabbitMQChannel,
		projectRepo:     projectRepo,
		courseRepo:      courseRepo,
		committeeRepo:   committeeRepo,
		programRepo:     programRepo,
	}
}

func (s *projectServiceImpl) PublishProjectMessageToElasticSearch(ctx context.Context, action string, projectId int) {
	go func() {
		projectMessage, err := s.projectRepo.GetProjectMessageByID(ctx, projectId)
		if err != nil {
			log.Printf("Failed to get project message")
		}

		if err = rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, action, projectMessage); err != nil {
			log.Printf("Failed to publish message to RabbitMQ for action %s: %v", action, err)
		}
	}()
}

func (s *projectServiceImpl) CreateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error) {

	projectMessage, err := s.projectRepo.CreateProjectWithFiles(ctx, project, projectResources, files)
	if err != nil {
		return nil, err
	}

	s.PublishProjectMessageToElasticSearch(ctx, "create", projectMessage.ID)

	return projectMessage, nil
}

func (s *projectServiceImpl) GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error) {
	project, err := s.projectRepo.GetProjectWithPDFByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) UpdateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error) {
	projectMessage, err := s.projectRepo.UpdateProjectWithFiles(ctx, project, projectResources, files)
	if err != nil {
		return nil, err
	}

	s.PublishProjectMessageToElasticSearch(ctx, "update", project.ID)

	return projectMessage, nil
}

func (s *projectServiceImpl) DeleteProject(ctx context.Context, id int) error {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.PublishProjectMessageToElasticSearch(ctx, "delete", project.ID)

	return nil
}
