package models

type ProjectRole struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	RoleName  string  `json:"role_name"`
	Program   Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID int     `json:"program_id"`
}
