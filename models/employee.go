package models

type Employee struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`                         // Employee ID (Primary Key)
	Prefix    string `json:"prefix"`                                                     // Prefix (e.g., Mr., Ms., Dr.)
	FirstName string `json:"first_name"`                                                 // First Name
	LastName  string `json:"last_name"`                                                  // Last Name
	Email     string `json:"email" gorm:"unique"`                                        // Unique Email
	RoleID    int    `json:"role_id" gorm:"not null"`                                    // Role ID (Foreign Key)
	Role      Role   `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL"` // Related Role (Foreign Key)
}
