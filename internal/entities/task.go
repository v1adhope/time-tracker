package entity

import "time"

type Task struct {
	ID         string    `json:"-"`
	CreatedAt  time.Time `json:"createdAt"`
	FinishedAt time.Time `json:"finishedAt"`
	// Description string    `json:"description"`
}
