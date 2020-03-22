package app

import (
	"time"

	"github.com/bobrovka/calendar/internal/models"
)

type App interface {
	ListDayEvents(date time.Time) ([]*models.Event, error)
	ListWeekEvents(date time.Time) ([]*models.Event, error)
	ListMonthEvents(date time.Time) ([]*models.Event, error)
	CreateNewEvent(newEvent *models.Event) (string, error)
	RemoveEvent(uuid string) error
	ChangeEvent(uuid string, newEvent *models.Event) error
}

// Calendar сущность, описывающая бизнес-логику сервиса
type Calendar struct {
	storage EventStorage
}

// NewCalendar создает новый инстанс приложения
func NewCalendar(storage EventStorage) (*Calendar, error) {
	return &Calendar{
		storage: storage,
	}, nil
}

// ListDayEvents вернет список событий на день
func (a *Calendar) ListDayEvents(date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(date, date.AddDate(0, 0, 1))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListWeekEvents вернет список событий на неделю
func (a *Calendar) ListWeekEvents(date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(date, date.AddDate(0, 0, 7))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListMonthEvents вернет список событий на месяц
func (a *Calendar) ListMonthEvents(date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(date, date.AddDate(0, 1, 0))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// CreateNewEvent добавит новое событие
func (a *Calendar) CreateNewEvent(newEvent *models.Event) (string, error) {
	// this should be like one transaction
	currentEvents, err := a.storage.ListEvents(newEvent.StartAt, newEvent.StartAt.AddDate(0, 0, 1))
	if err != nil {
		return "", err
	}

	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.StartAt.Add(newEvent.Duration)) {
		return "", ErrTimeBusy
	}

	return a.storage.CreateEvent(newEvent)
}

// RemoveEvent удалит событие
func (a *Calendar) RemoveEvent(uuid string) error {
	return a.storage.DeleteEvent(uuid)
}

// ChangeEvent изменит событие
func (a *Calendar) ChangeEvent(uuid string, newEvent *models.Event) error {
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
