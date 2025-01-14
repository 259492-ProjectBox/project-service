package dtos

type ResourceType struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	TypeName string `json:"type_name"`
}
