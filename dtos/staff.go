package dtos

type ProjectStaffMessage struct {
	ID          int         `json:"id"`
	Prefix      string      `json:"prefix"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Email       string      `json:"email"`
	ProgramID   int         `json:"program_id"`
	Program     Program     `json:"program"`
	ProjectRole ProjectRole `json:"project_role"`
}

type CreateStaffRequest struct {
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ProgramID int    `json:"program_id"`
}
type UpdateStaffRequest struct {
	ID        int    `json:"id"`
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ProgramID int    `json:"program_id"`
}

type StaffResponse struct {
	ID        int    `json:"id"`
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ProgramID int    `json:"program_id"`
}
