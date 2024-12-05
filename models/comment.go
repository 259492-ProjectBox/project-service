package models

import "time"

// Comment model
type Comment struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProjectID   int       `json:"project_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	Project     Project   `json:"project" gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}
