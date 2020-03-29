package models

import (
	"time"
)

// Event описывает событие
type Event struct {
	UUID         string
	Title        string
	StartAt      time.Time `db:"start_at"`
	Duration     time.Duration
	Description  string        `db:"descr"`
	User         string        `db:"user_name"`
	NotifyBefore time.Duration `db:"notify_before"`
}

func (e Event) String() string {
	return e.Title + " is starting at " + e.StartAt.Format("15:04:05")
}
