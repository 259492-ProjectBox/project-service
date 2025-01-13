package models

import "time"

type Calendar struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Program     Program   `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID   int       `json:"program_id"`
}

func (Calendar) TableName() string {
	return "calendar"
}
