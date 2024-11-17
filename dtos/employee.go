package dtos

type Employee struct {
	ID        int    `json:"id"`         // Employee ID (Primary Key)
	Prefix    string `json:"prefix"`     // Prefix (e.g., Mr., Ms., Dr.)
	FirstName string `json:"first_name"` // First Name
	LastName  string `json:"last_name"`  // Last Name
	Email     string `json:"email"`      // Unique Email
}
