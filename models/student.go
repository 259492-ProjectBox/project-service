package models

type Student struct {
	ID           int     `json:"id" gorm:"primaryKey;autoIncrement"`
	StudentID    string  `json:"student_id"`
	SecLab       string  `json:"sec_lab"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Semester     int     `json:"semester"`
	AcademicYear int     `json:"academic_year"`
	CourseID     int     `json:"course_id"`
	Course       Course  `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL"`
	ProgramID    int     `json:"program_id"`
	Program      Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
}
