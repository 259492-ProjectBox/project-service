package repositories

import (
	"context"
	"errors"
	"fmt"

	dto "github.com/project-box/dtos"
	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	repository[models.Project]
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]dto.ProjectWithDetails, error)
	GetByProjectNo(ctx context.Context, oldProjectNo string) (*models.Project, error)
	CreateProject(ctx context.Context, project *models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, id int, project *models.Project) (*models.Project, error)
}

// ProjectHandler implementation
type projectRepositoryImpl struct {
	db *gorm.DB
	*repositoryImpl[models.Project]
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepositoryImpl{
		db:             db,
		repositoryImpl: newRepository[models.Project](db),
	}
}

func (r *projectRepositoryImpl) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).Preload("Major").Preload("Course").Preload("Section").Preload("Advisor").Preload("Advisor.Role").Preload("Committees").Preload("Members").First(project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return project, nil
}
func (r *projectRepositoryImpl) GetProjectsByStudentId(ctx context.Context, studentId string) ([]dto.ProjectWithDetails, error) {
	projectsWithDetails := []dto.ProjectWithDetails{}

	var projectIds []int
	if err := r.db.WithContext(ctx).
		Model(&models.ProjectStudent{}).
		Where("project_students.student_id = ?", studentId).
		Pluck("project_students.project_id", &projectIds).Error; err != nil {
		return nil, err
	}

	if len(projectIds) == 0 {
		return nil, fmt.Errorf("no projects found for student ID %s", studentId)
	}

	for _, projectId := range projectIds {
		var project models.Project
		if err := r.db.WithContext(ctx).
			Where("id = ?", projectId).
			Preload("Advisor").      // Preload advisor
			Preload("Advisor.Role"). // Preload advisor
			Preload("Major").        // Preload major
			Preload("Course").       // Preload course
			Preload("Section").      // Preload section
			Find(&project).Error; err != nil {
			return nil, err
		}

		var employees []models.Employee
		if err := r.db.WithContext(ctx).
			Preload("Role").
			Joins("LEFT JOIN project_employees ON project_employees.employee_id = employees.id").
			Where("project_employees.project_id = ?", projectId).
			Find(&employees).Error; err != nil {
			return nil, err
		}

		var students []models.Student
		if err := r.db.WithContext(ctx).
			Preload("Major").
			Joins("LEFT JOIN project_students ON project_students.student_id = students.id").
			Where("project_students.project_id = ?", projectId).
			Find(&students).Error; err != nil {
			return nil, err
		}

		projectDetails := dto.ProjectWithDetails{
			Project:   project,
			Employees: employees,
			Students:  students,
		}

		projectsWithDetails = append(projectsWithDetails, projectDetails)
	}

	return projectsWithDetails, nil
}

func (r *projectRepositoryImpl) GetByProjectNo(ctx context.Context, projectNo string) (*models.Project, error) {
	project := &models.Project{}
	if err := r.db.WithContext(ctx).Where("project_no = ?", projectNo).First(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectRepositoryImpl) CreateProject(ctx context.Context, project *models.Project) (*models.Project, error) {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}
func (r *projectRepositoryImpl) UpdateProject(ctx context.Context, id int, project *models.Project) (*models.Project, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	existingProject := &models.Project{}
	if err := tx.First(existingProject, id).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("project not found")
	}
	if err := tx.Where("project_id = ?", id).Delete(&models.ProjectEmployee{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("project_id = ?", id).Delete(&models.ProjectStudent{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	project.ID = id
	if err := tx.Updates(project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	updatedProject := &models.Project{}
	if err := r.db.WithContext(ctx).
		Preload("Major").
		Preload("Course").
		Preload("Section").
		Preload("Advisor").
		Preload("Advisor.Role").
		Preload("Committees").
		Preload("Members").
		First(updatedProject, id).Error; err != nil {
		return nil, err
	}

	return updatedProject, nil
}
