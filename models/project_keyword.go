package models

type ProjectKeyword struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	ProjectID int     `json:"project_id"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	KeywordID int     `json:"keyword_id"`
	Keyword   Keyword `json:"keyword" gorm:"foreignKey:KeywordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
