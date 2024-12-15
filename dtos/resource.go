package dtos

type Resource struct {
	ID             int          `json:"id"`
	Title          string       `json:"title"`
	CreatedAt      string       `json:"created_at"`
	ResourceName   string       `json:"resource_name"`
	Path           string       `json:"path"`
	PDF            PDF          `json:"pdf"`
	ResourceTypeID int          `json:"resource_type_id"`
	ResourceType   ResourceType `json:"resource_type"`
}

type ProjectResource struct {
	ID       int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Resource Resource `json:"resource"`
}
