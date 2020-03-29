package app

import (
	"context"
	"time"

	"github.com/bobrovka/calendar/internal/models"
)

// EventStorage хранилище событий
type EventStorage interface {
	ListEvents(ctx context.Context, user string, from, to time.Time) ([]*models.Event, error)
	CreateEvent(ctx context.Context, event *models.Event) (string, error)
	UpdateEvent(ctx context.Context, id string, event *models.Event) error
	DeleteEvent(ctx context.Context, id string) error
}
