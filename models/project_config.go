package models

type ProjectConfig struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	Title     string  `json:"title"`
	IsActive  bool    `json:"is_active"`
	Program   Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID int     `json:"program_id"`
}
