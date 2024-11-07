package services

import (
	"context"
	"fmt"

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
func (s *projectServiceImpl) validateProjectNumber(projectNo string) error {
	if err := utils.IsValidProjectNumberFormat(projectNo); err != nil {
		return err
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
	if err := s.validateProjectNumber(project.ProjectNo); err != nil {
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

	err = rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, "create", project)
	if err != nil {
		return nil, err
	}

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

	project, err := s.projectRepo.Update(ctx, id, project)
	if err != nil {
		return nil, err
	}

	err = rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, "update", project)
	if err != nil {
		return nil, err
	}

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

	err = rabbitMQQueue.PublishMessageFromRabbitMQToElasticSearch(s.rabbitMQChannel, "delete", project)
	if err != nil {
		return err
	}

	return nil
}
