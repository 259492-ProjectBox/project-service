package models

type ProjectEmployee struct {
	ProjectID  int      `json:"project_id" gorm:"not null"`
	EmployeeID int      `json:"employee_id" gorm:"not null"`
	Project    Project  `json:"project" gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
	Employee   Employee `json:"employee" gorm:"foreignKey:EmployeeID;constraint:OnDelete:SET NULL"`
}
