package models

type ProjectNumberCounter struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Number       int    `gorm:"not null,default:1"`
	AcademicYear int    `gorm:"not null"`
	Semester     int    `gorm:"default:0"`
	CourseID     int    `gorm:"not null"`
	Course       Course `gorm:"foreignKey:CourseID;references:ID"`
}
