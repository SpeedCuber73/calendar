package app

// EventStorage хранилище событий
type EventStorage interface {
	ListEvents() ([]Event, error)
	CreateEvent(event *Event) error
	UpdateEvent(id string, event *Event) error
	DeleteEvent(id string) error
}
