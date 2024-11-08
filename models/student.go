package models

// Student model
type Student struct {
	ID        string `json:"id" gorm:"primaryKey"`                                         // Use `ID` as the primary key
	Prefix    string `json:"prefix"`                                                       // Optional field for title (e.g., Mr., Ms., etc.)
	FirstName string `json:"first_name"`                                                   // First name of the student
	LastName  string `json:"last_name"`                                                    // Last name of the student
	Email     string `json:"email" gorm:"unique"`                                          // Unique email address
	MajorID   int    `json:"major_id" gorm:"not null"`                                     // Major ID, not null
	Major     Major  `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:SET NULL"` // Major related to student
}
