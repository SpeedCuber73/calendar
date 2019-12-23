// Package simplestorage implements thread-unsafe storage
package simplestorage

import (
	app "github.com/SpeedCuber73/calendar/internal/calendar-app"
)

// SimpleStorage хранилище
type SimpleStorage struct {
	Events []app.Event
}

// CreateSimpleStorage конструктор
func CreateSimpleStorage() (*SimpleStorage, error) {
	return &SimpleStorage{
		Events: make([]app.Event, 0),
	}, nil
}

// ListEvents извлекает список событий
func (s *SimpleStorage) ListEvents() ([]app.Event, error) {
	toReturn := make([]app.Event, len(s.Events))
	copy(toReturn, s.Events)
	return toReturn, nil
}

// CreateEvent добавляет новое событие в хранилище
func (s *SimpleStorage) CreateEvent(event *app.Event) error {
	s.Events = append(s.Events, *event)
	return nil
}

// UpdateEvent обновляет информацию о событии
func (s *SimpleStorage) UpdateEvent(id string, renewEvent *app.Event) error {
	for i, event := range s.Events {
		if event.UUID == id {
			event = *renewEvent
			event.UUID = id
			s.Events[i] = event
			return nil
		}
	}
	return app.ErrNotFound
}

// DeleteEvent удаляет событие
func (s *SimpleStorage) DeleteEvent(id string) error {
	for i, event := range s.Events {
		if event.UUID == id {
			s.Events = append(s.Events[:i], s.Events[i+1:]...)
			return nil
		}
	}
	return app.ErrNotFound
}
