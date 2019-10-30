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
	ID          int
	Name        string
	StartAt     time.Time
	EndAt       time.Time
	Place       Place
	Description string
	CreatedAt   time.Time
}

func (e Event) String() string {
	return e.Name + " is starting at " + e.StartAt.String()
}
