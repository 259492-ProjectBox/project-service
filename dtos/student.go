package dtos

// Student model
type Student struct {
	ID        string `json:"id"`         // Use `ID` as the primary key
	FirstName string `json:"first_name"` // First name of the student
	LastName  string `json:"last_name"`  // Last name of the student
	Email     string `json:"email"`      // Unique email address
}
