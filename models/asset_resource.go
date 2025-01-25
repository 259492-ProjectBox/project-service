package models

type AssetResource struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Program     Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID   int     `json:"program_id"`
}
