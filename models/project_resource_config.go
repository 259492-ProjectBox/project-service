package models

type ProjectResourceConfig struct {
	ID             int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string       `json:"title"`
	ResourceTypeID int          `json:"resource_type_id"`
	ResourceType   ResourceType `json:"resource_type" gorm:"foreignKey:ResourceTypeID;constraint:OnDelete:CASCADE"`
	MajorID        int          `json:"major_id"`
	Major          Major        `json:"major" gorm:"foreignKey:MajorID;constraint:OnDelete:CASCADE"`
}
