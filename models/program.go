package models

type Program struct {
	ID            int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProgramNameTH string `json:"program_name_th" gorm:"unique"`
	ProgramNameEN string `json:"program_name_en" gorm:"unique"`
}
