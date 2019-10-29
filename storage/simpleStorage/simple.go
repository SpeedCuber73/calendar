package simplestorage

import (
	"github.com/SpeedCuber73/calendar/storage"
)

// SimpleStorage хранилище
type SimpleStorage struct {
	Events []storage.Event
}

// CreateSimpleStorage конструктор
func CreateSimpleStorage() (*SimpleStorage, error) {
	return &SimpleStorage{
		Events: make([]storage.Event, 0),
	}, nil
}

// ListEvents извлекает список событий
func (s *SimpleStorage) ListEvents() ([]storage.Event, error) {
	toReturn := make([]storage.Event, 0, len(s.Events))
	copy(toReturn, s.Events)
	return toReturn, nil
}

// CreateEvent добавляет новое событие в хранилище
func (s *SimpleStorage) CreateEvent(event *storage.Event) error {
	s.Events = append(s.Events, *event)
	return nil
}

// UpdateEvent обновляет информацию о событии
func (s *SimpleStorage) UpdateEvent(id int, renewEvent *Event) error {
	for i, event := range s.Events {
		if event.Id == id {
			event = renewEvent
			event.Id = id
			s.Events[i] = event
			return nil
		}
	}
	return storage.ErrNotFound
}

// DeleteEvent удаляет событие
func (s *SimpleStorage) DeleteEvent(id int) error {
	for i, event := range s.Events {
		if event.Id == id {
			s.Events = append(s.Events[:i], s.Events[i+1:])
			return nil
		}
	}
	return storage.ErrNotFound
}
