package app

import (
	"time"
)

// Place описывает место
type Place struct {
	Country     string
	City        string
	Street      string
	Building    int
	Description string
}

// Event описывает событие
type Event struct {
	UUID        string
	Title       string
	StartAt     time.Time
	EndAt       time.Time
	Place       Place
	Description string
	CreatedAt   time.Time
}

func (e Event) String() string {
	return e.Title + " is starting at " + e.StartAt.String()
}
