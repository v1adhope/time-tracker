package entities

type Task struct {
	ID         string `json:"id" example:"1ef4e803-1af7-6a50-85b2-77ed6f34a8cf"`
	CreatedAt  string `json:"createdAt" example:"2024-01-16 09:08:25"`
	FinishedAt string `json:"finishedAt,omitempty" example:"2024-01-16 16:10:00"`
	UserID     string `json:"userId" example:"1ef4e803-1aed-62e0-8d59-c8cfd7561759"`
}

type TaskSummary struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	FinishedAt  string `json:"finishedAt,omitempty"`
	SummaryTime string `json:"summaryTime,omitempty"`
}

type TaskSort struct {
	StartTime string
	EndTime   string
}
