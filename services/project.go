package services

import (
	"context"
	"fmt"
	"log"

	dto "github.com/project-box/dtos"
	"github.com/project-box/models"
	rabbitMQQueue "github.com/project-box/queues/rabbitmq"
	"github.com/project-box/repositories"
	"github.com/project-box/utils"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type ProjectService interface {
	ValidateProject(ctx context.Context, project *models.Project) error
	CreateProject(ctx context.Context, project *models.Project) (*models.Project, error)
	GetProjectById(ctx context.Context, id int) (*models.Project, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]dto.ProjectWithDetails, error)
	UpdateProject(ctx context.Context, id int, project *models.Project) (*models.Project, error)
	DeleteProject(ctx context.Context, id int) error
}

type projectServiceImpl struct {
	rabbitMQChannel *rabbitmq.Channel
	projectRepo     repositories.ProjectRepository
	committeeRepo   repositories.EmployeeRepository
	majorRepo       repositories.MajorRepository
	sectionRepo     repositories.SectionRepository
}

func NewProjectService(
	rabbitMQChannel *rabbitmq.Channel,
	projectRepo repositories.ProjectRepository,
	committeeRepo repositories.EmployeeRepository,
	majorRepo repositories.MajorRepository,
	sectionRepo repositories.SectionRepository,
) ProjectService {
	return &projectServiceImpl{
		rabbitMQChannel: rabbitMQChannel,
		projectRepo:     projectRepo,
		committeeRepo:   committeeRepo,
		majorRepo:       majorRepo,
		sectionRepo:     sectionRepo,
	}
}

// Helper function to publish CRUD operations to Elasticsearch via RabbitMQ
func (s *projectServiceImpl) publishToElasticSearch(ctx context.Context, action string, project *models.Project) {
	go func() {
		projectMessage := utils.SanitizeProjectMessage(project)
		err := rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, action, projectMessage)
		if err != nil {
			log.Printf("Failed to publish message to RabbitMQ for action %s: %v", action, err)
		}
	}()
}

func (s *projectServiceImpl) createProjectNumber(project *models.Project) *models.Project {
	projectNumber := utils.FormatProjectNumber(*project.SectionID, project.Semester, project.AcademicYear)
	project.ProjectNo = projectNumber
	return project
}

func (s *projectServiceImpl) validateCourse(ctx context.Context, courseID int, sectionID *int, semester int) error {
	section, err := s.sectionRepo.GetByCourseAndSectionAndSemester(ctx, courseID, sectionID, semester)
	if err != nil {
		return fmt.Errorf("invalid course and section combination")
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

func (s *projectServiceImpl) validateOldProjectNumber(ctx context.Context, oldProjectNo *string) error {
	if oldProjectNo == nil {
		return nil
	}

	if err := utils.IsValidProjectNumberFormat(*oldProjectNo); err != nil {
		return err
	}

	_, err := s.projectRepo.GetByProjectNo(ctx, *oldProjectNo)
	if err != nil {
		return fmt.Errorf("project number not exists: %w", err)
	}

	return nil
}
func (s *projectServiceImpl) ValidateProject(ctx context.Context, project *models.Project) error {
	if err := s.validateCourse(ctx, project.CourseID, project.SectionID, project.Semester); err != nil {
		return err
	}
	if err := s.validateMajor(ctx, project.MajorID); err != nil {
		return err
	}
	if err := s.validateAdvisor(ctx, project.AdvisorID); err != nil {
		return err
	}
	if err := s.validateOldProjectNumber(ctx, project.OldProjectNo); err != nil {
		return err
	}
	return nil
}

func (s *projectServiceImpl) CreateProject(ctx context.Context, project *models.Project) (*models.Project, error) {
	if err := s.ValidateProject(ctx, project); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.Create(ctx, s.createProjectNumber(project))
	if err != nil {
		return nil, err
	}

	s.publishToElasticSearch(ctx, "create", project)

	return project, nil
}

func (s *projectServiceImpl) GetProjectById(ctx context.Context, id int) (*models.Project, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) GetProjectsByStudentId(ctx context.Context, studentId string) ([]dto.ProjectWithDetails, error) {
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

	// Instead of deleting and re-creating, update the project directly
	project, err := s.projectRepo.UpdateProject(ctx, id, project)
	if err != nil {
		return nil, err
	}

	s.publishToElasticSearch(ctx, "update", project)

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

	s.publishToElasticSearch(ctx, "delete", project)

	return nil
}
