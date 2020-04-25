package scheduler

import (
	"context"

	"github.com/bobrovka/calendar/internal/models"
)

// EventStorage хранилище событий
type EventStorage interface {
	PopNotifications(ctx context.Context) ([]*models.Event, error)
}
