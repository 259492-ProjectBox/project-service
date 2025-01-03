package models

type ProjectConfig struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	IsActive bool   `json:"is_active"`
	Major    Major  `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`
	MajorID  int    `json:"major_id"`
}
