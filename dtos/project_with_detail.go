package dtos

import "github.com/project-box/models"

type ProjectWithDetails struct {
	Project   models.Project    `json:"project"`
	Employees []models.Employee `json:"employees"`
	Students  []models.Student  `json:"students"`
}
