package storage

import (
	"time"
)

type Place struct {
	Country     string
	City        string
	Street      string
	Building    int
	Description string
}

type Event struct {
	Id          int
	Name        string
	StartAt     time.Time
	EndAt       time.Time
	Place       Place
	Description string
	CreatedAt   time.Time
}
