package models

type DetailedResource struct {
	Project         Project         `gorm:"embedded"`
	ProjectResource ProjectResource `gorm:"embedded"`
}
