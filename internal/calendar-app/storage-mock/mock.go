package mock

import (
	"time"

	"github.com/bobrovka/calendar/internal/models"
	"github.com/stretchr/testify/mock"
)

// StorageMock мок хранилища
type StorageMock struct {
	mock.Mock
}

// ListEvents мокирует метод
func (m *StorageMock) ListEvents(from, to time.Time) ([]*models.Event, error) {
	args := m.Called(from, to)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	return args.Get(0).([]*models.Event), err
}

// CreateEvent мокирует метод
func (m *StorageMock) CreateEvent(event *models.Event) (string, error) {
	args := m.Called(event)
	err := args.Error(1)
	if err != nil {
		return "", err
	}

	return args.String(0), err
}

// UpdateEvent мокирует метод
func (m *StorageMock) UpdateEvent(id string, event *models.Event) error {
	args := m.Called(id, event)
	return args.Error(0)
}

// DeleteEvent мокирует метод
func (m *StorageMock) DeleteEvent(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
