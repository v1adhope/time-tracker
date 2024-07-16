package entities

type Task struct {
	ID         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	FinishedAt string `json:"finishedAt,omitempty"`
	UserID     string `json:"userId"`
	// Description string    `json:"description"`
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
