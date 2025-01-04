package dtos

// Major model
type Major struct {
	ID        int    `json:"id"`
	MajorName string `json:"major_name"`
}

type CreateMajorRequest struct {
	MajorName string `json:"major_name" binding:"required"`
}
