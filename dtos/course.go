package dtos

// Course model
type Course struct {
	ID         int    `json:"id"`
	CourseNo   string `json:"course_no"`
	CourseName string `json:"course_name"`
	MajorID    int    `json:"major_id"`
	Major      Major  `json:"major" `
}
