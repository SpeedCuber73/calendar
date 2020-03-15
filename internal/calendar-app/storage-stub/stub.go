package stub

import (
	"time"

	"github.com/bobrovka/calendar/internal/models"
)

// StorageStub is dummy storage
type StorageStub struct{}

// ListEvents stub for method
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

// CreateEvent stub for method
func (s *StorageStub) CreateEvent(_ *models.Event) (string, error) {
	return "1", nil
}

// UpdateEvent stub for method
func (s *StorageStub) UpdateEvent(_ string, _ *models.Event) error {
	return nil
}

// DeleteEvent stub for method
func (s *StorageStub) DeleteEvent(_ string) error {
	return nil
}
