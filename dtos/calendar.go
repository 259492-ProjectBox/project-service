package dtos

type CreateCalendarRequest struct {
	MajorID     int    `json:"major_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date" `
}

type CalendarResponse struct {
	ID          int    `json:"id"`
	StartDate   string `json:"start_date" `
	EndDate     string `json:"end_date" `
	Title       string `json:"title"`
	Description string `json:"description"`
	Major       string `json:"major_name"`
}

type UpdateCalendarRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date" `
	MajorID     int    `json:"major_id"`
}
