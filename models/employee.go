package models

type Employee struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix    string    `json:"prefix"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"unique"`
	MajorID   int       `json:"major_id"`
	Major     Major     `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:SET NULL"`
	Projects  []Project `json:"projects" gorm:"many2many:project_employees;constraint:OnDelete:CASCADE;"`
}
