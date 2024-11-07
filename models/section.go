package models

// Section model
type Section struct {
	ID            int    `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID      int    `json:"course_id" gorm:"not null"`
	SectionNumber string `json:"section_number"`
	Semester      int    `json:"semester"`
	Course        Course `json:"course" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}
