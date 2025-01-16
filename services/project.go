package services

import (
	"context"
	"fmt"
	"log"

	"github.com/project-box/dtos"
	"github.com/project-box/models"
	rabbitMQQueue "github.com/project-box/queues/rabbitmq"
	"github.com/project-box/repositories"
	"github.com/project-box/utils"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type ProjectService interface {
	PublishProjectMessageToElasticSearch(ctx context.Context, action string, projectId int)
	ValidateProject(ctx context.Context, project *models.ProjectRequest) error
	GetProjectById(ctx context.Context, id int) (*dtos.ProjectData, error)
	GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error)
	// GetProjectByAdvisorId(ctx context.Context, advisorId int) ([]dtos.ProjectData, error)
	CreateProjectWithFiles(ctx context.Context, req *dtos.CreateProjectRequest) (*dtos.ProjectData, error)
	UpdateProjectWithFiles(ctx context.Context, req *dtos.UpdateProjectRequest) (*dtos.ProjectData, error)
	DeleteProject(ctx context.Context, id int) error
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

func (s *projectServiceImpl) createProjectNumber(project *models.ProjectRequest) (*models.ProjectRequest, error) {
	nextProjectNumber, err := s.projectNumberCounterRepo.GetNextProjectNumber(project.AcademicYear, project.Semester, project.CourseID)
	if err != nil {
		return nil, err
	}
	projectNumber := utils.FormatProjectID(project.Semester, project.AcademicYear, nextProjectNumber)
	project.ProjectNo = projectNumber
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

func (s *projectServiceImpl) ValidateProject(ctx context.Context, project *models.ProjectRequest) error {
	if err := s.validateCourse(ctx, project.CourseID, project.Semester); err != nil {
		return err
	}
	if err := s.validateProgram(ctx, project.ProgramID); err != nil {
		return err
	}
	return nil
}
func (s *projectServiceImpl) CreateProjectWithFiles(ctx context.Context, req *dtos.CreateProjectRequest) (*dtos.ProjectData, error) {
	project := req.Project
	files := req.Files
	titles := req.Titles
	urls := req.Urls

	if err := s.ValidateProject(ctx, project); err != nil {
		return nil, err
	}

	project, err := s.createProjectNumber(project)
	if err != nil {
		return nil, err
	}
	projectMessage, err := s.projectRepo.CreateProjectWithFiles(ctx, project, files, titles, urls)
	if err != nil {
		return nil, err
	}
	s.PublishProjectMessageToElasticSearch(ctx, "create", projectMessage.ID)

	return projectMessage, nil
}

func (s *projectServiceImpl) GetProjectById(ctx context.Context, id int) (*dtos.ProjectData, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectServiceImpl) GetProjectWithPDFByID(ctx context.Context, id int) (*dtos.ProjectData, error) {
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

// func (s *projectServiceImpl) GetProjectsByAdvisorId(ctx context.Context, advisorId int) ([]dtos.ProjectData, error) {
// 	project, err := s.projectRepo.GetProjectsByAdvisorId(ctx, advisorId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return project, nil
// }

func (s *projectServiceImpl) UpdateProjectWithFiles(ctx context.Context, req *dtos.UpdateProjectRequest) (*dtos.ProjectData, error) {
	project := req.Project
	files := req.Files
	titles := req.Titles
	urls := req.Urls
	if err := s.ValidateProject(ctx, project); err != nil {
		return nil, err
	}

	projectMessage, err := s.projectRepo.UpdateProjectWithFiles(ctx, project, files, titles, urls)
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
