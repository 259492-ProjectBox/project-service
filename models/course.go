package models

type Course struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseNo   string `json:"course_no" gorm:"unique"`
	CourseName string `json:"course_name" gorm:"unique"`
	MajorID    int    `json:"major_id" gorm:"not null"`
	Major      Major  `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`
}
