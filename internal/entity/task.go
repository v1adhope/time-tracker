package entity

import "time"

type Task struct {
	ID        string    `json:"-"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	// Description string    `json:"description"`
}
