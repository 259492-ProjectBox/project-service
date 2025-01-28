package models

type Course struct {
	ID         int     `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseNo   string  `json:"course_no"`
	CourseName string  `json:"course_name"`
	ProgramID  int     `json:"program_id"`
	Program    Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
}
