package models

import "time"

type Calendar struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	EventDate   time.Time `json:"event_date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Major       Major     `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`
	MajorID     int       `json:"major_id"`
}
