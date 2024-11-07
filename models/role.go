package models

// Role model
type Role struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleName string `json:"role_name" gorm:"unique"`
}
