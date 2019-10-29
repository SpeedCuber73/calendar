package app

// EventStorage хранилище событий
type EventStorage interface {
	ListEvents() ([]Event, error)
	CreateEvent(event *Event) error
	UpdateEvent(id int, event *Event) error
	DeleteEvent(id int) error
}
