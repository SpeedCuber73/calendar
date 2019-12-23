package app

import "time"

// App сущность, описывающая бизнес-логику сервиса
type App struct {
	storage EventStorage
}

// NewApp создает новый инстанс приложения
func NewApp(storage EventStorage) (*App, error) {
	return &App{
		storage: storage,
	}, nil
}

// ListAllEvents вернет список событий
func (a *App) ListAllEvents() []Event {
	events, err := a.storage.ListEvents()
	if err != nil {
		return nil
	}
	return events
}

// AddNewEvent добавит новое событие
func (a *App) AddNewEvent(newEvent *Event) error {
	// this should be like one transaction
	currentEvents, err := a.storage.ListEvents()
	if err != nil {
		return err
	}
	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.EndAt) {
		return ErrTimeBusy
	}
	return a.storage.CreateEvent(newEvent)
}

// RemoveEvent удалит событие
func (a *App) RemoveEvent(id string) error {
	return a.storage.DeleteEvent(id)
}

// ChangeEvent изменит событие
func (a *App) ChangeEvent(id string, newEvent *Event) error {
	// this should be like one transaction
	currentEvents, err := a.storage.ListEvents()
	if err != nil {
		return err
	}

	for i, event := range currentEvents {
		if event.UUID == id {
			currentEvents = append(currentEvents[:i], currentEvents[i+1:]...)
			break
		}
	}

	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.EndAt) {
		return ErrTimeBusy
	}
	return a.storage.UpdateEvent(id, newEvent)
}

func hasFreeTime(existingEvents []Event, start, end time.Time) bool {
	for _, event := range existingEvents {
		if (event.StartAt.Before(start) || event.StartAt.Equal(start)) && event.EndAt.After(start) {
			return false
		}
		if event.StartAt.After(start) && event.StartAt.Before(end) {
			return false
		}
	}
	return true
}
