package models

type ProjectEmployeeType struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	TypeName string `json:"typeName"`
}
