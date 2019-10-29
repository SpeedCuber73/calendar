package app

import "errors"

var (
	// ErrNotFound объект не найден
	ErrNotFound = errors.New("object not found")
)

// EventStorage хранилище событий
type EventStorage interface {
	ListEvents() ([]Event, error)
	CreateEvent(event *Event) error
	UpdateEvent(id int, event *Event) error
	DeleteEvent(id int) error
}
