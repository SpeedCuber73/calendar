package app

import (
	"time"

	"github.com/bobrovka/calendar/internal/models"
)

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

// ListAllEvents вернет список событий на день
func (a *App) ListDayEvents(date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(date, date.AddDate(0, 0, 1))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListAllEvents вернет список событий на неделю
func (a *App) ListWeekEvents(date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(date, date.AddDate(0, 0, 7))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListAllEvents вернет список событий на месяц
func (a *App) ListMonthEvents(date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(date, date.AddDate(0, 1, 0))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// CreateNewEvent добавит новое событие
func (a *App) CreateNewEvent(newEvent *models.Event) error {
	// this should be like one transaction
	currentEvents, err := a.storage.ListEvents(newEvent.StartAt, newEvent.StartAt.AddDate(0, 0, 1))
	if err != nil {
		return err
	}

	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.StartAt.Add(newEvent.Duration)) {
		return ErrTimeBusy
	}

	return a.storage.CreateEvent(newEvent)
}

// RemoveEvent удалит событие
func (a *App) RemoveEvent(uuid string) error {
	return a.storage.DeleteEvent(uuid)
}

// ChangeEvent изменит событие
func (a *App) ChangeEvent(uuid string, newEvent *models.Event) error {
	// get events on this day
	currentEvents, err := a.storage.ListEvents(newEvent.StartAt, newEvent.StartAt.AddDate(0, 0, 1))
	if err != nil {
		return err
	}

	var found bool
	// delete an event that is being modified
	for i, event := range currentEvents {
		if event.UUID == uuid {
			currentEvents = append(currentEvents[:i], currentEvents[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return ErrNotFound
	}

	// if no free time - abort changing
	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.StartAt.Add(newEvent.Duration)) {
		return ErrTimeBusy
	}

	return a.storage.UpdateEvent(uuid, newEvent)
}

func hasFreeTime(existingEvents []*models.Event, start, end time.Time) bool {
	for _, event := range existingEvents {
		eventEndAt := event.StartAt.Add(event.Duration)
		if (event.StartAt.Before(start) || event.StartAt.Equal(start)) && eventEndAt.After(start) {
			return false
		}
		if event.StartAt.After(start) && event.StartAt.Before(end) {
			return false
		}
	}

	return true
}
