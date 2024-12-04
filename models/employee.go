package models

type Employee struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	RoleID    int    `json:"role_id" gorm:"not null"`
	Role      Role   `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL"`
	MajorID   int    `json:"major_id" gorm:"not null"`
	Major     Major  `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:SET NULL"`
}
