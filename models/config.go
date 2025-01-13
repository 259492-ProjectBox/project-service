package models

type Config struct {
	ConfigName string  `json:"config_name"`
	Value      string  `json:"value"`
	Program    Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID  int     `json:"program_id"`
}
