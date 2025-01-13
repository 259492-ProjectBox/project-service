package dtos

// Course model
type Course struct {
	ID         int     `json:"id"`
	CourseNo   string  `json:"course_no"`
	CourseName string  `json:"course_name"`
	ProgramID  int     `json:"program_id"`
	Program    Program `json:"program"`
}
