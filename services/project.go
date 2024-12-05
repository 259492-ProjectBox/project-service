package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/project-box/models"
	rabbitMQQueue "github.com/project-box/queues/rabbitmq"
	"github.com/project-box/repositories"
	"github.com/project-box/utils"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type ProjectService interface {
	ValidateProject(ctx context.Context, project *models.Project) error
	CreateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string) (*models.Project, error)
	GetProjectById(ctx context.Context, id int) (*models.Project, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error)
	UpdateProject(ctx context.Context, id int, project *models.Project) (*models.Project, error)
	DeleteProject(ctx context.Context, id int) error
}

type projectServiceImpl struct {
	rabbitMQChannel          *rabbitmq.Channel
	projectRepo              repositories.ProjectRepository
	committeeRepo            repositories.EmployeeRepository
	majorRepo                repositories.MajorRepository
	courseRepo               repositories.CourseRepository
	projectNumberCounterRepo repositories.ProjectNumberCounterRepository
}

func NewProjectService(
	rabbitMQChannel *rabbitmq.Channel,
	projectRepo repositories.ProjectRepository,
	committeeRepo repositories.EmployeeRepository,
	majorRepo repositories.MajorRepository,
	courseRepo repositories.CourseRepository,
	projectNumberCounterRepo repositories.ProjectNumberCounterRepository,
) ProjectService {
	return &projectServiceImpl{
		rabbitMQChannel:          rabbitMQChannel,
		projectRepo:              projectRepo,
		courseRepo:               courseRepo,
		committeeRepo:            committeeRepo,
		majorRepo:                majorRepo,
		projectNumberCounterRepo: projectNumberCounterRepo,
	}
}

func (s *projectServiceImpl) publishProjectMessageToElasticSearch(ctx context.Context, action string, project *models.Project) {
	go func() {
		project, err := s.GetProjectById(ctx, project.ID)
		if err != nil {
			log.Printf("Failed to fetch project by ID: %v", err)
			return
		}

		if project == nil {
			log.Printf("Project with ID is nil")
			return
		}

		projectMessage := utils.SanitizeProjectMessage(project)
		fmt.Printf("%+v\n", projectMessage)
		if err = rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, action, projectMessage); err != nil {
			log.Printf("Failed to publish message to RabbitMQ for action %s: %v", action, err)
		}
	}()
}

func (s *projectServiceImpl) createProjectNumber(project *models.Project) (*models.Project, error) {
	nextProjectNumber, err := s.projectNumberCounterRepo.GetNextProjectNumber(project.AcademicYear, project.Semester, project.CourseID)
	if err != nil {
		return nil, err
	}
	projectID := utils.FormatProjectID(project.Semester, project.AcademicYear, nextProjectNumber)
	project.ProjectNo = projectID
	return project, nil
}

func (s *projectServiceImpl) validateCourse(ctx context.Context, courseID int, semester int) error {
	section, err := s.courseRepo.GetByCourseAndSemester(ctx, courseID, semester)
	if err != nil {
		return fmt.Errorf("invalid course and semester combination")
	}
	if section == nil {
		return fmt.Errorf("section not found")
	}
	return nil
}

func (s *projectServiceImpl) validateMajor(ctx context.Context, majorID int) error {
	major, err := s.majorRepo.Get(ctx, majorID)
	if err != nil || major == nil {
		return fmt.Errorf("major not found")
	}
	return nil
}

func (s *projectServiceImpl) validateAdvisor(ctx context.Context, advisorID int) error {
	advisor, err := s.committeeRepo.Get(ctx, advisorID)
	if err != nil || advisor == nil {
		return fmt.Errorf("advisor not found")
	}
	return nil
}

func (s *projectServiceImpl) ValidateProject(ctx context.Context, project *models.Project) error {
	if err := s.validateCourse(ctx, project.CourseID, project.Semester); err != nil {
		return err
	}
	if err := s.validateMajor(ctx, project.MajorID); err != nil {
		return err
	}
	return nil
}

func (s *projectServiceImpl) CreateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string) (*models.Project, error) {
	if err := s.ValidateProject(ctx, project); err != nil {
		return nil, err
	}

	project, err := s.createProjectNumber(project)
	if err != nil {
		return nil, err
	}

	project, err = s.projectRepo.CreateProjectWithFiles(ctx, project, files, titles)
	if err != nil {
		return nil, err
	}

	s.publishProjectMessageToElasticSearch(ctx, "create", project)

	return project, nil
}

func (s *projectServiceImpl) GetProjectById(ctx context.Context, id int) (*models.Project, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error) {
	project, err := s.projectRepo.GetProjectsByStudentId(ctx, studentId)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) UpdateProject(ctx context.Context, id int, project *models.Project) (*models.Project, error) {
	if err := s.ValidateProject(ctx, project); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.UpdateProject(ctx, id, project)
	if err != nil {
		return nil, err
	}

	s.publishProjectMessageToElasticSearch(ctx, "update", project)

	return project, nil
}

func (s *projectServiceImpl) DeleteProject(ctx context.Context, id int) error {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.publishProjectMessageToElasticSearch(ctx, "delete", project)

	return nil
}
