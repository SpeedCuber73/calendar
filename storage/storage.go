package storage

import "errors"

var (
	ErrNotFound = errors.New("object not found")
)

type EventStorage interface {
	ListEvents() ([]Event, error)
	CreateEvent(event *Event) error
	UpdateEvent(id int, event *Event) error
	DeleteEvent(id int) error
}
