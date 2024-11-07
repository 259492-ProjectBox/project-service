package models

import "time"

// Employee model
type Employee struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	EmployeeName string    `json:"employee_name"`
	Email        string    `json:"email" gorm:"unique"`
	RoleID       int       `json:"role_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	Role         Role      `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL"`
}
