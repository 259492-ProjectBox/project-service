package dtos

type ProjectConfigResponse struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	MajorID  int    `json:"major_id"`
	IsActive bool   `json:"is_active"`
}

type ProjectConfigUpsertRequest struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	MajorID  int    `json:"major_id"`
	IsActive bool   `json:"is_active"`
}

type InsertProjectConfigRequest struct {
	Title    string `json:"title"`
	MajorID  int    `json:"major_id"`
	IsActive bool   `json:"is_active"`
}
