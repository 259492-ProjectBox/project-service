package dtos

type Program struct {
	ID            int    `json:"id"`
	ProgramNameTH string `json:"program_name_th"`
	ProgramNameEN string `json:"program_name_en"`
}

type CreateProgramRequest struct {
	ProgramNameTH string `json:"program_name_th"`
	ProgramNameEN string `json:"program_name_en"`
}
