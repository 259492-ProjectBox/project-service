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
	PublishProjectMessageToElasticSearch(ctx context.Context, action string, project *models.Project)
	ValidateProject(ctx context.Context, project *models.Project) error
	CreateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string) (*models.Project, error)
	UpdateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string) (*models.Project, error)
	GetProjectById(ctx context.Context, id int) (*models.Project, error)
	GetProjectWithPDFByID(ctx context.Context, id int) (*models.Project, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error)
	DeleteProject(ctx context.Context, id int) error
	GetProjectsByAdvisorIdService(ctx context.Context, advisorId int) ([]models.Project, error)
}

type projectServiceImpl struct {
	rabbitMQChannel          *rabbitmq.Channel
	projectRepo              repositories.ProjectRepository
	committeeRepo            repositories.StaffRepository
	programRepo              repositories.ProgramRepository
	courseRepo               repositories.CourseRepository
	projectNumberCounterRepo repositories.ProjectNumberCounterRepository
}

func NewProjectService(
	rabbitMQChannel *rabbitmq.Channel,
	projectRepo repositories.ProjectRepository,
	committeeRepo repositories.StaffRepository,
	programRepo repositories.ProgramRepository,
	courseRepo repositories.CourseRepository,
	projectNumberCounterRepo repositories.ProjectNumberCounterRepository,
) ProjectService {
	return &projectServiceImpl{
		rabbitMQChannel:          rabbitMQChannel,
		projectRepo:              projectRepo,
		courseRepo:               courseRepo,
		committeeRepo:            committeeRepo,
		programRepo:              programRepo,
		projectNumberCounterRepo: projectNumberCounterRepo,
	}
}

func (s *projectServiceImpl) PublishProjectMessageToElasticSearch(ctx context.Context, action string, project *models.Project) {
	go func() {
		project, err := s.GetProjectWithPDFByID(ctx, project.ID)
		if err != nil {
			log.Printf("Failed to fetch project by ID: %v", err)
			return
		}

		if project == nil {
			log.Printf("Project with ID is nil")
			return
		}

		projectMessage := utils.SanitizeProjectMessage(project)
		fmt.Printf("%v", projectMessage)
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

func (s *projectServiceImpl) validateProgram(ctx context.Context, programId int) error {
	program, err := s.programRepo.Get(ctx, programId)
	if err != nil || program == nil {
		return fmt.Errorf("program not found")
	}
	return nil
}

// func (s *projectServiceImpl) validateAdvisor(ctx context.Context, advisorID int) error {
// 	advisor, err := s.committeeRepo.Get(ctx, advisorID)
// 	if err != nil || advisor == nil {
// 		return fmt.Errorf("advisor not found")
// 	}
// 	return nil
// }

func (s *projectServiceImpl) ValidateProject(ctx context.Context, project *models.Project) error {
	if err := s.validateCourse(ctx, project.CourseID, project.Semester); err != nil {
		return err
	}
	if err := s.validateProgram(ctx, project.ProgramID); err != nil {
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

	s.PublishProjectMessageToElasticSearch(ctx, "create", project)

	return project, nil
}

func (s *projectServiceImpl) GetProjectById(ctx context.Context, id int) (*models.Project, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}
func (s *projectServiceImpl) GetProjectWithPDFByID(ctx context.Context, id int) (*models.Project, error) {
	project, err := s.projectRepo.GetProjectWithPDFByID(ctx, id)
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

func (s *projectServiceImpl) UpdateProjectWithFiles(ctx context.Context, project *models.Project, files []*multipart.FileHeader, titles []string) (*models.Project, error) {
	if err := s.ValidateProject(ctx, project); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.UpdateProjectWithFiles(ctx, project, files, titles)
	if err != nil {
		return nil, err
	}

	s.PublishProjectMessageToElasticSearch(ctx, "update", project)

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

	s.PublishProjectMessageToElasticSearch(ctx, "delete", project)

	return nil
}

// get project relate to advisor from advisor id
func (s *projectServiceImpl) GetProjectsByAdvisorIdService(ctx context.Context, advisorId int) ([]models.Project, error) {
	project, err := s.projectRepo.GetProjectsByAdvisorId(ctx, advisorId)
	if err != nil {
		return nil, err
	}

	return project, nil
}
