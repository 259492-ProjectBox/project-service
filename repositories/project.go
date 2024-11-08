package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/project-box/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	repository[models.Project]
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error)
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
func (r *projectRepositoryImpl) GetProjectsByStudentId(ctx context.Context, studentId string) ([]models.Project, error) {
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

	var projects []models.Project
	if err := r.db.WithContext(ctx).
		Where("id IN ?", projectIds).
		Preload("Advisor").
		Preload("Advisor.Role").
		Preload("Major").
		Preload("Course").
		Preload("Section").
		Find(&projects).Error; err != nil {
		return nil, err
	}

	var projectEmployees []models.ProjectEmployee
	if err := r.db.WithContext(ctx).
		Model(&models.ProjectEmployee{}).
		Preload("Employee").
		Preload("Employee.Role").
		Where("project_employees.project_id IN ?", projectIds).
		Find(&projectEmployees).Error; err != nil {
		return nil, err
	}

	employeesMap := make(map[int][]models.Employee)
	for _, projectEmployee := range projectEmployees {
		projectID := projectEmployee.ProjectID
		employee := projectEmployee.Employee
		employeesMap[projectID] = append(employeesMap[projectID], employee)
	}

	var projectStudents []models.ProjectStudent
	if err := r.db.WithContext(ctx).
		Model(&models.ProjectStudent{}).
		Preload("Student").
		Preload("Student.Major").
		Where("project_students.project_id IN ?", projectIds).
		Find(&projectStudents).Error; err != nil {
		return nil, err
	}

	studentsMap := make(map[int][]models.Student)
	for _, projectStudent := range projectStudents {
		projectID := projectStudent.ProjectID
		studentsMap[projectID] = append(studentsMap[projectID], projectStudent.Student)
	}

	for i := range projects {
		project := &projects[i]
		project.Committees = employeesMap[project.ID]
		project.Members = studentsMap[project.ID]
	}

	return projects, nil
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
