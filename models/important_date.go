package models

import "time"

// ImportantDate model
type ImportantDate struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	MajorID     int       `json:"major_id" gorm:"not null"`
	EventDate   time.Time `json:"event_date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Major       Major     `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`
}
