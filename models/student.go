package models

// Student model
type Student struct {
	ID          string `json:"id" gorm:"primaryKey"`
	StudentName string `json:"student_name"`
	Email       string `json:"email" gorm:"unique"`
	MajorID     int    `json:"major_id" gorm:"not null"`
	Major       Major  `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:SET NULL"`
}
