package dtos

import (
	"time"
)

type CreateCalendarRequest struct {
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date" `
	Title       string `json:"title"`
	Description string `json:"description"`
	MajorID     int    `json:"major_id"`
}

type CalendarResponse struct {
	ID          int       `json:"id"`
	StartDate   time.Time `json:"start_date" gorm:"default:CURRENT_TIMESTAMP"`
	EndDate     time.Time `json:"end_date" gorm:"default:CURRENT_TIMESTAMP"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Major       string    `json:"major_name"`
}
