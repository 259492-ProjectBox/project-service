package models

type ProjectNumberCounter struct {
	ID           int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Number       int     `json:"number" gorm:"default:1"`
	AcademicYear int     `json:"academic_year"`
	Semester     int     `gorm:"default:0"`
	ProgramID    *int    `json:"program_id"`
	Program      Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
}
