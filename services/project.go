package services

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	rabbitMQQueue "github.com/project-box/queues/rabbitmq"
	"github.com/project-box/repositories"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ProjectService interface {
	PublishProjectMessageToElasticSearch(ctx context.Context, action string, projectId int) error
	GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error)
	CheckDuplicateProjectByTitleAndSemester(ctx context.Context, titleTH, titleEN string, academicYear, semester int) (bool, error)
	CreateProjects(ctx context.Context, project []models.ProjectRequest) error
	CreateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error)
	UpdateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error)
	DeleteProject(ctx context.Context, id int) error
}

type projectServiceImpl struct {
	rabbitMQChannel *rabbitmq.Channel
	projectRepo     repositories.ProjectRepository
	committeeRepo   repositories.StaffRepository
	programRepo     repositories.ProgramRepository
}

func NewProjectService(
	rabbitMQChannel *rabbitmq.Channel,
	projectRepo repositories.ProjectRepository,
	committeeRepo repositories.StaffRepository,
	programRepo repositories.ProgramRepository,
) ProjectService {
	return &projectServiceImpl{
		rabbitMQChannel: rabbitMQChannel,
		projectRepo:     projectRepo,
		committeeRepo:   committeeRepo,
		programRepo:     programRepo,
	}
}

func (s *projectServiceImpl) CreateProjects(ctx context.Context, projects []models.ProjectRequest) error {
	projectMessages, err := s.projectRepo.CreateProjects(ctx, projects)
	if err != nil {
		return err
	}
	for _, projectMessage := range projectMessages {
		err := s.PublishProjectMessageToElasticSearch(ctx, "create", projectMessage.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *projectServiceImpl) getProjectMessage(ctx context.Context, projectId int) (*dtos.ProjectData, error) {
	projectMessage, err := s.projectRepo.GetProjectMessageByID(ctx, projectId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if projectMessage == nil {
		projectMessage = &dtos.ProjectData{
			ID: projectId,
		}
	}

	return projectMessage, nil
}

func (s *projectServiceImpl) publishMessage(ctx context.Context, action string, projectMessage *dtos.ProjectData) error {
	if err := rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, action, projectMessage); err != nil {
		return err
	}

	return nil
}

func (s *projectServiceImpl) PublishProjectMessageToElasticSearch(ctx context.Context, action string, projectId int) error {
	projectMessage, err := s.getProjectMessage(ctx, projectId)
	if err != nil {
		return err
	}

	return s.publishMessage(ctx, action, projectMessage)
}

func (s *projectServiceImpl) CreateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error) {
	projectMessage, err := s.projectRepo.CreateProjectWithFiles(ctx, nil, project, projectResources, files)
	if err != nil {
		return nil, err
	}

	err = s.PublishProjectMessageToElasticSearch(ctx, "create", projectMessage.ID)
	if err != nil {
		return nil, err
	}

	return projectMessage, nil
}

func (s *projectServiceImpl) GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error) {
	project, err := s.projectRepo.GetProjectWithPDFByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) CheckDuplicateProjectByTitleAndSemester(ctx context.Context, titleTH, titleEN string, academicYear, semester int) (bool, error) {
	isDuplicate, err := s.projectRepo.CheckDuplicateProjectByTitleAndSemester(ctx, titleTH, titleEN, academicYear, semester)
	if err != nil {
		return true, err
	}

	return isDuplicate, nil
}

func (s *projectServiceImpl) UpdateProjectWithFiles(ctx context.Context, project *models.ProjectRequest, projectResources []*models.ProjectResource, files []*multipart.FileHeader) (*dtos.ProjectData, error) {
	projectMessage, err := s.projectRepo.UpdateProjectWithFiles(ctx, nil, project, projectResources, files)
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

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return err
	}
	if err := s.PublishProjectMessageToElasticSearch(ctx, "delete", id); err != nil {
		return err
	}
	return nil
}
