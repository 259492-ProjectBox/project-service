package models

type ProjectStaff struct {
	ID            int         `json:"id" gorm:"primaryKey"`
	ProjectID     int         `json:"project_id"`
	Project       Project     `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProjectRoleID int         `json:"project_role_id"`
	ProjectRole   ProjectRole `json:"project_role" gorm:"foreignKey:ProjectRoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StaffID       int         `json:"staff_id"`
	Staff         Staff       `json:"staff" gorm:"foreignKey:StaffID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
