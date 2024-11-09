package models

// Course model
type Course struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseNo   string `json:"course_no"`
	CourseName string `json:"course_name" gorm:"unique"`
}
