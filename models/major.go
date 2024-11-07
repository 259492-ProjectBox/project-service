package models

// Major model
type Major struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	MajorName string `json:"major_name" gorm:"unique"`
}
