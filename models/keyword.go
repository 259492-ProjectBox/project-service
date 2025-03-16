package models

type Keyword struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	Keyword   string  `json:"keyword"`
	ProgramID int     `json:"program_id"`
	Program   Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
