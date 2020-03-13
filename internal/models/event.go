package models

import (
	"time"
)

// Event описывает событие
type Event struct {
	UUID         string
	Title        string
	StartAt      time.Time
	Duration     time.Duration
	Description  string
	User         string
	NotifyBefore time.Duration
}

func (e Event) String() string {
	return e.Title + " is starting at " + e.StartAt.Format("15:04:05")
}
