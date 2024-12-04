package models

type ProjectEmployeeType struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	TypeName string `json:"type_name" gorm:"unique"`
}
