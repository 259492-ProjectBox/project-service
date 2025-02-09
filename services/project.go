package services

import (
	"context"
	"mime/multipart"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	rabbitMQQueue "github.com/project-box/queues/rabbitmq"
	"github.com/project-box/repositories"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type ProjectService interface {
	PublishProjectMessageToElasticSearch(ctx context.Context, action string, projectId int) error
	GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error)
	CreateProjects(ctx context.Context, project []models.Project) error
	CreateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error)
	CreateProjectsWithFilesTx(ctx context.Context, projects []models.ProjectRequest, projectResources [][]*models.ProjectResource, files [][]*multipart.FileHeader) ([]*dtos.ProjectData, error)
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

func (s *projectServiceImpl) CreateProjects(ctx context.Context, project []models.Project) error {
	if err := s.projectRepo.CreateMany(ctx, project); err != nil {
		return err
	}
	return nil
}

func (s *projectServiceImpl) PublishProjectMessageToElasticSearch(ctx context.Context, action string, projectId int) error {
	projectMessage, err := s.projectRepo.GetProjectMessageByID(ctx, projectId)
	if err != nil {
		return err
	}

	if err = rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, action, projectMessage); err != nil {
		return err
	}

	return nil
}

func (s *projectServiceImpl) CreateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error) {
	projectMessage, err := s.projectRepo.CreateProjectWithFiles(ctx, project, projectResources, files)
	if err != nil {
		return nil, err
	}

	err = s.PublishProjectMessageToElasticSearch(ctx, "create", projectMessage.ID)
	if err != nil {
		return nil, err
	}

	return projectMessage, nil
}

func (s *projectServiceImpl) CreateProjectsWithFilesTx(ctx context.Context, projects []models.ProjectRequest, projectResources [][]*models.ProjectResource, files [][]*multipart.FileHeader) ([]*dtos.ProjectData, error) {
	var projectMessages []*dtos.ProjectData
	for i, project := range projects {
		var resources []*models.ProjectResource
		var fileHeaders []*multipart.FileHeader

		if projectResources != nil && i < len(projectResources) {
			resources = projectResources[i]
		}

		if files != nil && i < len(files) {
			fileHeaders = files[i]
		}

		projectMessage, err := s.projectRepo.CreateProjectWithFiles(ctx, &project, resources, fileHeaders)
		if err != nil {
			return nil, err
		}

		projectMessages = append(projectMessages, projectMessage)
	}

	for _, projectMessage := range projectMessages {
		s.PublishProjectMessageToElasticSearch(ctx, "create", projectMessage.ID)
	}

	return projectMessages, nil
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

	err = s.PublishProjectMessageToElasticSearch(ctx, "update", project.ID)
	if err != nil {
		return nil, err
	}

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

	err = s.PublishProjectMessageToElasticSearch(ctx, "delete", project.ID)
	if err != nil {
		return err
	}

	return nil
}
