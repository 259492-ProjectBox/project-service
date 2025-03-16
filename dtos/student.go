package dtos

type Student struct {
	ID           int     `json:"id"`         // Use `ID` as the primary key
	StudentID    string  `json:"student_id"` // Student ID
	SecLab       string  `json:"sec_lab"`
	FirstName    string  `json:"first_name"`    // First name of the student
	LastName     string  `json:"last_name"`     // Last name of the student
	Email        *string `json:"email"`         // Unique email address
	Semester     int     `json:"semester"`      // Semester
	AcademicYear int     `json:"academic_year"` // Academic year
	ProgramID    int     `json:"program_id"`    // Program ID
	Program      Program `json:"program"`
}
