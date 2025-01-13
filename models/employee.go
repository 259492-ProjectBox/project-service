package models

type Employee struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix    string    `json:"prefix"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"unique"`
	Program   Program   `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID int       `json:"program_id"`
	Projects  []Project `json:"projects" gorm:"many2many:project_employees;constraint:OnDelete:CASCADE;"`
}
