package models

type ProjectEmployee struct {
	ID                    int                 `json:"id" gorm:"primaryKey"`
	ProjectID             int                 `json:"project_id"`
	EmployeeID            int                 `json:"employee_id"`
	Project               Project             `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Employee              Employee            `json:"employee" gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProjectEmployeeTypeID int                 `json:"project_employee_type_id"`
	ProjectEmployeeType   ProjectEmployeeType `json:"project_employee_type" gorm:"foreignKey:ProjectEmployeeTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
