package models

type ProjectNumberCounter struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Number       int    `json:"number" gorm:"default:1"`
	AcademicYear int    `json:"academic_year"`
	Semester     int    `gorm:"default:0"`
	CourseID     int    `json:"course_id"`
	Course       Course `json:"course" gorm:"foreignKey:CourseID;references:ID"`
}
