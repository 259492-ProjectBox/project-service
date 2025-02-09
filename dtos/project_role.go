package dtos

type ProjectRole struct {
	ID         int     `json:"id"`
	RoleNameTH string  `json:"role_name_th"`
	RoleNameEN string  `json:"role_name_en"`
	Program    Program `json:"program"`
	ProgramID  int     `json:"program_id"`
}
