package app

import (
	"time"

	"github.com/bobrovka/calendar/internal/models"
)

// EventStorage хранилище событий
type EventStorage interface {
	ListEvents(from, to time.Time) ([]*models.Event, error)
	CreateEvent(event *models.Event) (string, error)
	UpdateEvent(id string, event *models.Event) error
	DeleteEvent(id string) error
}
