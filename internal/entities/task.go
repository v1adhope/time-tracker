package entities

import "time"

type Task struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	FinishedAt time.Time `json:"finishedAt"`
	UserID     string    `json:"userId"`
	// Description string    `json:"description"`
}
