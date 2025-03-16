package dtos

type Keyword struct {
	ID        int     `json:"id"`
	Keyword   string  `json:"keyword"`
	ProgramID int     `json:"program_id"`
	Program   Program `json:"program"`
}
