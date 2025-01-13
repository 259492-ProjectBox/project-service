package dtos

type Program struct {
	ID          int    `json:"id"`
	ProgramName string `json:"program_name"`
}

type CreateProgramRequest struct {
	ProgramName string `json:"program_name"`
}
