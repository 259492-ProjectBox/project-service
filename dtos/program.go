package dtos

type Program struct {
	ID            int    `json:"id"`
	ProgramNameTH string `json:"program_name_th" gorm:"unique"`
	ProgramNameEN string `json:"program_name_en" gorm:"unique"`
}

type CreateProgramRequest struct {
	ProgramNameTH string `json:"program_name_th" gorm:"unique"`
	ProgramNameEN string `json:"program_name_en" gorm:"unique"`
}
