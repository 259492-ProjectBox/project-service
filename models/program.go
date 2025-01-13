package models

// Program model
type Program struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProgramName string `json:"program_name" gorm:"unique"`
}
