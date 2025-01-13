package dtos

type Employee struct {
	ID        int     `json:"id"`         // Employee ID (Primary Key)
	Prefix    string  `json:"prefix"`     // Prefix (e.g., Mr., Ms., Dr.)
	FirstName string  `json:"first_name"` // First Name
	LastName  string  `json:"last_name"`  // Last Name
	Email     string  `json:"email"`      // Unique Email
	ProgramID int     `json:"program_id"`
	Program   Program `json:"program"`
}

type CreateEmployeeRequest struct {
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	MajorID   int    `json:"major_id"`
}
type UpdateEmployeeRequest struct {
	ID        int    `json:"id"`
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	MajorID   int    `json:"major_id"`
}

type EmployeeResponse struct {
	ID        int    `json:"id"`
	Prefix    string `json:"prefix"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	MajorID   int    `json:"major_id"`
}
