package stub

import (
	"time"

	"github.com/bobrovka/calendar/internal/models"
)

type StorageStub struct{}

func (s *StorageStub) ListEvents(_, _ time.Time) ([]*models.Event, error) {
	return []*models.Event{
		&models.Event{
			UUID:         "1",
			Title:        "title-1",
			StartAt:      time.Now(),
			Duration:     2 * time.Hour,
			Description:  "awesome meeting",
			User:         "Kira",
			NotifyBefore: 3 * time.Hour,
		},
	}, nil
}

func (s *StorageStub) CreateEvent(_ *models.Event) error {
	return nil
}

func (s *StorageStub) UpdateEvent(_ string, _ *models.Event) error {
	return nil
}

func (s *StorageStub) DeleteEvent(_ string) error {
	return nil
}
