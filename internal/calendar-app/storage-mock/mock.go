package mock

import (
	"time"

	"github.com/SpeedCuber73/calendar/internal/models"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (m *StorageMock) ListEvents(from, to time.Time) ([]*models.Event, error) {
	args := m.Called(from, to)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	return args.Get(0).([]*models.Event), err
}

func (m *StorageMock) CreateEvent(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *StorageMock) UpdateEvent(id string, event *models.Event) error {
	args := m.Called(id, event)
	return args.Error(0)
}

func (m *StorageMock) DeleteEvent(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
