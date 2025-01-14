package models

// Student model
type Student struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
}
