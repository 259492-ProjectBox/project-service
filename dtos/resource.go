package dtos

type Resource struct {
	ID             int          `json:"id"`
	Title          string       `json:"title"`
	URL            string       `json:"url"`
	CreatedAt      string       `json:"created_at"`
	PDF            PDF          `json:"pdf"`
	ResourceTypeID int          `json:"resource_type_id"`
	ResourceType   ResourceType `json:"resource_type"`
}

type ProjectResource struct {
	ID       int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Resource Resource `json:"resource"`
}
