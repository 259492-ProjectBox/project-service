package dtos

type Section struct {
	ID            int    `json:"id"`
	SectionNumber string `json:"section_number"`
	Semester      int    `json:"semester"`
}
