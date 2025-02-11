package models

type ProjectRole struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	RoleNameTH string  `json:"role_name_th"`
	RoleNameEN string  `json:"role_name_en"`
	Program    Program `json:"program" gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE"`
	ProgramID  int     `json:"program_id"`
}
