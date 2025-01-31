package dtos

type ProjectStaffMessage struct {
	ID          int         `json:"id"`
	PrefixTH    string      `json:"prefix_th"`
	PrefixEN    string      `json:"prefix_en"`
	FirstNameTH string      `json:"first_name_th"`
	LastNameTH  string      `json:"last_name_th"`
	FirstNameEN string      `json:"first_name_en"`
	LastNameEN  string      `json:"last_name_en"`
	Email       string      `json:"email"`
	ProgramID   int         `json:"program_id"`
	Program     Program     `json:"program"`
	ProjectRole ProjectRole `json:"project_role"`
}

type CreateStaffRequest struct {
	PrefixTH    string `json:"prefix_th"`
	PrefixEN    string `json:"prefix_en"`
	FirstNameTH string `json:"first_name_th"`
	LastNameTH  string `json:"last_name_th"`
	FirstNameEN string `json:"first_name_en"`
	LastNameEN  string `json:"last_name_en"`
	Email       string `json:"email"`
	ProgramID   int    `json:"program_id"`
}

type UpdateStaffRequest struct {
	ID          int    `json:"id"`
	PrefixTH    string `json:"prefix_th"`
	PrefixEN    string `json:"prefix_en"`
	FirstNameTH string `json:"first_name_th"`
	LastNameTH  string `json:"last_name_th"`
	FirstNameEN string `json:"first_name_en"`
	LastNameEN  string `json:"last_name_en"`
	Email       string `json:"email"`
	ProgramID   int    `json:"program_id"`
}

type StaffResponse struct {
	ID          int    `json:"id"`
	PrefixTH    string `json:"prefix_th"`
	PrefixEN    string `json:"prefix_en"`
	FirstNameTH string `json:"first_name_th"`
	LastNameTH  string `json:"last_name_th"`
	FirstNameEN string `json:"first_name_en"`
	LastNameEN  string `json:"last_name_en"`
	Email       string `json:"email"`
	ProgramID   int    `json:"program_id"`
}
