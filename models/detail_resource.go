package models

type DetailedResource struct {
	Project         Project         `gorm:"embedded"`
	Resource        Resource        `gorm:"embedded"`
	ProjectResource ProjectResource `gorm:"embedded"`
	AssetResource   AssetResource   `gorm:"embedded"`
}
