package models

type ProjectEmployee struct {
	ID         int      `json:"id" gorm:"primaryKey"`
	ProjectID  int      `json:"project_id"`
	EmployeeID int      `json:"employee_id"`
	Project    Project  `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Employee   Employee `json:"employee" gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
